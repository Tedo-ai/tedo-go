package tedo

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStorageService_PutObjectEscapesSlashSeparatedKeys(t *testing.T) {
	const key = "inbox/2026/04/message.eml"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s, want PUT", r.Method)
		}
		if got, want := r.URL.EscapedPath(), "/storage/v1/buckets/bucket-1/objects/inbox%2F2026%2F04%2Fmessage.eml"; got != want {
			t.Fatalf("escaped path = %q, want %q", got, want)
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !bytes.Equal(data, []byte("hello")) {
			t.Fatalf("body = %q, want %q", data, "hello")
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":           "obj-1",
			"bucket_id":    "bucket-1",
			"key":          key,
			"size":         5,
			"content_type": "message/rfc822",
			"hash":         "hash-1",
			"created_at":   "2026-01-01T00:00:00Z",
		})
	}))
	defer server.Close()

	client := NewClient("tedo_live_test").WithBaseURL(server.URL).WithHTTPClient(server.Client())
	obj, err := client.Storage.PutObject(context.Background(), "bucket-1", key, strings.NewReader("hello"), "message/rfc822")
	if err != nil {
		t.Fatalf("PutObject failed: %v", err)
	}
	if obj.Key != key {
		t.Fatalf("returned key = %q, want %q", obj.Key, key)
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
