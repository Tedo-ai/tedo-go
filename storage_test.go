package tedo

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewClient_DefaultHTTPClientHasNoTimeout(t *testing.T) {
	client := NewClient("tedo_live_test")
	if client.httpClient.Timeout != 0 {
		t.Fatalf("default timeout = %s, want 0", client.httpClient.Timeout)
	}
}

func TestStorageService_ObjectKeyPathEscapesPerSegment(t *testing.T) {
	const key = "assets/browser-X.js"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Regression for tedo-core PR #256: object keys are escaped per segment
		// so literal slashes remain path separators and never become %2F.
		if got, want := r.URL.EscapedPath(), "/storage/v1/buckets/bucket-1/objects/assets/browser-X.js"; got != want {
			t.Fatalf("escaped path = %q, want %q", got, want)
		}

		switch r.Method {
		case http.MethodPut:
			data, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			if !bytes.Equal(data, []byte("console.log('x');")) {
				t.Fatalf("body = %q, want %q", data, "console.log('x');")
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":           "obj-1",
				"bucket_id":    "bucket-1",
				"key":          key,
				"size":         17,
				"content_type": "application/javascript",
				"hash":         "hash-1",
				"created_at":   "2026-01-01T00:00:00Z",
			})
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/javascript")
			_, _ = w.Write([]byte("console.log('x');"))
		case http.MethodHead:
			w.Header().Set("Content-Type", "application/javascript")
			w.Header().Set("Content-Length", "17")
			w.Header().Set("ETag", `"hash-1"`)
			w.WriteHeader(http.StatusOK)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("method = %s, want PUT/GET/HEAD/DELETE", r.Method)
		}
	}))
	defer server.Close()

	client := NewClient("tedo_live_test").WithBaseURL(server.URL).WithHTTPClient(server.Client())
	obj, err := client.Storage.PutObject(context.Background(), "bucket-1", key, strings.NewReader("console.log('x');"), "application/javascript")
	if err != nil {
		t.Fatalf("PutObject failed: %v", err)
	}
	if obj.Key != key {
		t.Fatalf("returned key = %q, want %q", obj.Key, key)
	}

	body, contentType, err := client.Storage.GetObject(context.Background(), "bucket-1", key)
	if err != nil {
		t.Fatalf("GetObject failed: %v", err)
	}
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("read downloaded body: %v", err)
	}
	if string(data) != "console.log('x');" {
		t.Fatalf("downloaded body = %q", data)
	}
	if contentType != "application/javascript" {
		t.Fatalf("download content-type = %q", contentType)
	}

	head, err := client.Storage.HeadObject(context.Background(), "bucket-1", key)
	if err != nil {
		t.Fatalf("HeadObject failed: %v", err)
	}
	if head.Key != key {
		t.Fatalf("head key = %q, want %q", head.Key, key)
	}
	if head.Size != 17 {
		t.Fatalf("head size = %d, want 17", head.Size)
	}

	if err := client.Storage.DeleteObject(context.Background(), "bucket-1", key); err != nil {
		t.Fatalf("DeleteObject failed: %v", err)
	}
}

func TestStorageService_PutObjectWithOptionsSendsHashHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Header.Get("X-Content-Sha256"), "abc123"; got != want {
			t.Fatalf("X-Content-Sha256 = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":           "obj-1",
			"bucket_id":    "bucket-1",
			"key":          "message.eml",
			"size":         5,
			"content_type": "message/rfc822",
			"hash":         "abc123",
			"created_at":   "2026-01-01T00:00:00Z",
		})
	}))
	defer server.Close()

	client := NewClient("tedo_live_test").WithBaseURL(server.URL).WithHTTPClient(server.Client())
	_, err := client.Storage.PutObjectWithOptions(context.Background(), "bucket-1", "message.eml", strings.NewReader("hello"), &PutObjectOptions{
		ContentType:   "message/rfc822",
		ContentSHA256: "abc123",
	})
	if err != nil {
		t.Fatalf("PutObjectWithOptions failed: %v", err)
	}
}

func TestStorageService_ListObjectsEncodesPrefixQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if got, want := r.URL.RawQuery, "limit=100&prefix=inbox%2F2026%2F04"; got != want {
			t.Fatalf("raw query = %q, want %q", got, want)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"objects": []map[string]any{},
			"limit":   100,
			"offset":  0,
		})
	}))
	defer server.Close()

	client := NewClient("tedo_live_test").WithBaseURL(server.URL).WithHTTPClient(server.Client())
	_, err := client.Storage.ListObjects(context.Background(), "bucket-1", &ListObjectsParams{
		Prefix: "inbox/2026/04",
		Limit:  100,
	})
	if err != nil {
		t.Fatalf("ListObjects failed: %v", err)
	}
}

func TestStorageService_HeadObject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodHead {
			t.Fatalf("method = %s, want HEAD", r.Method)
		}
		w.Header().Set("Content-Type", "message/rfc822")
		w.Header().Set("Content-Length", "42")
		w.Header().Set("ETag", `"hash-1"`)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("tedo_live_test").WithBaseURL(server.URL).WithHTTPClient(server.Client())
	obj, err := client.Storage.HeadObject(context.Background(), "bucket-1", "inbox/2026/04/message.eml")
	if err != nil {
		t.Fatalf("HeadObject failed: %v", err)
	}
	if obj.Key != "inbox/2026/04/message.eml" {
		t.Fatalf("key = %q", obj.Key)
	}
	if obj.Size != 42 {
		t.Fatalf("size = %d, want 42", obj.Size)
	}
	if obj.ContentType != "message/rfc822" {
		t.Fatalf("content-type = %q", obj.ContentType)
	}
	if obj.Hash != `"hash-1"` {
		t.Fatalf("hash = %q", obj.Hash)
	}
}

func TestStorageService_PutObjectRetriesTransientFailure(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := atomic.AddInt32(&attempts, 1)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !bytes.Equal(data, []byte("hello")) {
			t.Fatalf("body = %q, want %q", data, "hello")
		}

		if current == 1 {
			w.Header().Set("Retry-After", "0")
			http.Error(w, "try again", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":           "obj-1",
			"bucket_id":    "bucket-1",
			"key":          "message.eml",
			"size":         5,
			"content_type": "message/rfc822",
			"hash":         "hash-1",
			"created_at":   "2026-01-01T00:00:00Z",
		})
	}))
	defer server.Close()

	client := NewClient("tedo_live_test").
		WithBaseURL(server.URL).
		WithHTTPClient(server.Client()).
		WithRetryConfig(RetryConfig{MaxRetries: 1, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond})

	obj, err := client.Storage.PutObject(context.Background(), "bucket-1", "message.eml", strings.NewReader("hello"), "message/rfc822")
	if err != nil {
		t.Fatalf("PutObject failed: %v", err)
	}
	if obj.Key != "message.eml" {
		t.Fatalf("key = %q, want %q", obj.Key, "message.eml")
	}
	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Fatalf("attempts = %d, want 2", got)
	}
}

func TestStorageService_HeadObjectRetriesTransientFailure(t *testing.T) {
	var attempts int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := atomic.AddInt32(&attempts, 1)
		if current == 1 {
			http.Error(w, "temporary failure", http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "message/rfc822")
		w.Header().Set("Content-Length", "42")
		w.Header().Set("ETag", `"hash-1"`)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("tedo_live_test").
		WithBaseURL(server.URL).
		WithHTTPClient(server.Client()).
		WithRetryConfig(RetryConfig{MaxRetries: 1, InitialBackoff: time.Millisecond, MaxBackoff: time.Millisecond})

	obj, err := client.Storage.HeadObject(context.Background(), "bucket-1", "message.eml")
	if err != nil {
		t.Fatalf("HeadObject failed: %v", err)
	}
	if obj.Size != 42 {
		t.Fatalf("size = %d, want 42", obj.Size)
	}
	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Fatalf("attempts = %d, want 2", got)
	}
}
