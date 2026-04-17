package tedo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// StorageService handles object storage API calls.
type StorageService struct {
	client *Client
}

// --- Bucket types ---

// Bucket represents a storage bucket.
type Bucket struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	Visibility  string `json:"visibility"`
	ObjectCount int64  `json:"object_count"`
	TotalSize   int64  `json:"total_size"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateBucketParams are the parameters for creating a bucket.
type CreateBucketParams struct {
	Name       string `json:"name"`
	Visibility string `json:"visibility,omitempty"` // "private" (default) or "public-read"
}

// --- Object types ---

// Object represents a stored object.
type Object struct {
	ID          string `json:"id"`
	BucketID    string `json:"bucket_id"`
	Key         string `json:"key"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	Hash        string `json:"hash"`
	CreatedAt   string `json:"created_at"`
}

// ObjectList is the response from listing objects.
type ObjectList struct {
	Objects []*Object `json:"objects"`
	Limit   int       `json:"limit"`
	Offset  int       `json:"offset"`
}

// ListObjectsParams are the parameters for listing objects.
type ListObjectsParams struct {
	Prefix string
	Limit  int
	Offset int
}

// PutObjectOptions configure upload behavior.
type PutObjectOptions struct {
	ContentType   string
	ContentSHA256 string
}

// PresignParams are the parameters for creating a pre-signed URL.
type PresignParams struct {
	Key    string `json:"key"`
	Expiry int    `json:"expiry,omitempty"` // seconds, default 3600
}

// PresignResponse is the response from creating a pre-signed URL.
type PresignResponse struct {
	URL string `json:"url"`
}

// StorageUsage represents storage usage stats.
type StorageUsage struct {
	WorkspaceID  string `json:"workspace_id"`
	TotalSize    int64  `json:"total_size"`
	TotalObjects int64  `json:"total_objects"`
}

// --- Bucket methods ---

// ListBuckets lists all buckets.
func (s *StorageService) ListBuckets(ctx context.Context) ([]*Bucket, error) {
	var buckets []*Bucket
	if err := s.doJSON(ctx, http.MethodGet, "/storage/v1/buckets", nil, nil, &buckets); err != nil {
		return nil, err
	}
	return buckets, nil
}

// CreateBucket creates a new bucket.
func (s *StorageService) CreateBucket(ctx context.Context, params *CreateBucketParams) (*Bucket, error) {
	var bucket Bucket
	if err := s.client.request(ctx, http.MethodPost, "/storage/v1/buckets", params, &bucket); err != nil {
		return nil, err
	}
	return &bucket, nil
}

// GetBucket gets a bucket by ID.
func (s *StorageService) GetBucket(ctx context.Context, bucketID string) (*Bucket, error) {
	var bucket Bucket
	if err := s.doJSON(ctx, http.MethodGet, fmt.Sprintf("/storage/v1/buckets/%s", bucketID), nil, nil, &bucket); err != nil {
		return nil, err
	}
	return &bucket, nil
}

// DeleteBucket deletes an empty bucket.
func (s *StorageService) DeleteBucket(ctx context.Context, bucketID string) error {
	return s.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/storage/v1/buckets/%s", bucketID), nil, nil, nil)
}

// --- Object methods ---

// ListObjects lists objects in a bucket.
func (s *StorageService) ListObjects(ctx context.Context, bucketID string, params *ListObjectsParams) (*ObjectList, error) {
	path := fmt.Sprintf("/storage/v1/buckets/%s/objects", bucketID)
	if params != nil {
		query := url.Values{}
		if params.Prefix != "" {
			query.Set("prefix", params.Prefix)
		}
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			query.Set("offset", strconv.Itoa(params.Offset))
		}
		if encoded := query.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}

	var list ObjectList
	if err := s.doJSON(ctx, http.MethodGet, path, nil, nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// PutObject uploads an object. Returns metadata about the stored object.
func (s *StorageService) PutObject(ctx context.Context, bucketID, key string, body io.Reader, contentType string) (*Object, error) {
	return s.PutObjectWithOptions(ctx, bucketID, key, body, &PutObjectOptions{ContentType: contentType})
}

// PutObjectWithOptions uploads an object with additional upload options.
func (s *StorageService) PutObjectWithOptions(ctx context.Context, bucketID, key string, body io.Reader, opts *PutObjectOptions) (*Object, error) {
	payload, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("read request body: %w", err)
	}

	path := fmt.Sprintf("/storage/v1/buckets/%s/objects/%s", bucketID, url.PathEscape(key))
	contentType := "application/octet-stream"
	if opts != nil && opts.ContentType != "" {
		contentType = opts.ContentType
	}
	headers := map[string]string{
		"Content-Type": contentType,
	}
	if opts != nil && opts.ContentSHA256 != "" {
		headers["X-Content-Sha256"] = opts.ContentSHA256
	}

	var obj Object
	if err := s.doJSON(ctx, http.MethodPut, path, func() io.Reader {
		return bytes.NewReader(payload)
	}, headers, &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// HeadObject fetches object metadata without downloading the body.
func (s *StorageService) HeadObject(ctx context.Context, bucketID, key string) (*Object, error) {
	path := fmt.Sprintf("/storage/v1/buckets/%s/objects/%s", bucketID, url.PathEscape(key))
	resp, err := s.do(ctx, http.MethodHead, path, nil, nil)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, parseError(resp.StatusCode, []byte(http.StatusText(resp.StatusCode)))
	}

	var size int64
	if cl := resp.Header.Get("Content-Length"); cl != "" {
		size, _ = strconv.ParseInt(cl, 10, 64)
	}

	return &Object{
		BucketID:    bucketID,
		Key:         key,
		Size:        size,
		ContentType: resp.Header.Get("Content-Type"),
		Hash:        resp.Header.Get("ETag"),
	}, nil
}

// GetObject downloads an object. The caller must close the returned ReadCloser.
func (s *StorageService) GetObject(ctx context.Context, bucketID, key string) (io.ReadCloser, string, error) {
	path := fmt.Sprintf("/storage/v1/buckets/%s/objects/%s", bucketID, url.PathEscape(key))
	resp, err := s.do(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, "", parseError(resp.StatusCode, body)
	}

	return resp.Body, resp.Header.Get("Content-Type"), nil
}

// DeleteObject deletes an object.
func (s *StorageService) DeleteObject(ctx context.Context, bucketID, key string) error {
	return s.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/storage/v1/buckets/%s/objects/%s", bucketID, url.PathEscape(key)), nil, nil, nil)
}

// --- Pre-signed URLs ---

// PresignURL creates a pre-signed URL for temporary public access.
func (s *StorageService) PresignURL(ctx context.Context, bucketID string, params *PresignParams) (*PresignResponse, error) {
	var resp PresignResponse
	if err := s.client.request(ctx, http.MethodPost, fmt.Sprintf("/storage/v1/buckets/%s/presign", bucketID), params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Usage ---

// GetUsage returns storage usage statistics.
func (s *StorageService) GetUsage(ctx context.Context) (*StorageUsage, error) {
	var usage StorageUsage
	if err := s.doJSON(ctx, http.MethodGet, "/storage/v1/usage", nil, nil, &usage); err != nil {
		return nil, err
	}
	return &usage, nil
}

func (s *StorageService) doJSON(ctx context.Context, method, path string, bodyFactory func() io.Reader, headers map[string]string, result any) error {
	resp, err := s.do(ctx, method, path, bodyFactory, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return parseError(resp.StatusCode, respBody)
	}
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

func (s *StorageService) do(ctx context.Context, method, path string, bodyFactory func() io.Reader, headers map[string]string) (*http.Response, error) {
	cfg := normalizeRetryConfig(s.client.retry)
	for attempt := 0; ; attempt++ {
		req, err := s.newRequest(ctx, method, path, bodyFactory, headers)
		if err != nil {
			return nil, err
		}

		resp, err := s.client.httpClient.Do(req)
		if !shouldRetryStorageRequest(ctx, cfg, attempt, resp, err) {
			if err != nil {
				return nil, fmt.Errorf("do request: %w", err)
			}
			return resp, nil
		}

		delay := storageRetryDelay(cfg, attempt, resp)
		if resp != nil && resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
		if err := sleepWithContext(ctx, delay); err != nil {
			return nil, err
		}
	}
}

func (s *StorageService) newRequest(ctx context.Context, method, path string, bodyFactory func() io.Reader, headers map[string]string) (*http.Request, error) {
	var body io.Reader
	if bodyFactory != nil {
		body = bodyFactory()
	}
	req, err := http.NewRequestWithContext(ctx, method, s.client.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.client.apiKey)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return req, nil
}

func shouldRetryStorageRequest(ctx context.Context, cfg RetryConfig, attempt int, resp *http.Response, err error) bool {
	if attempt >= cfg.MaxRetries {
		return false
	}
	if ctx.Err() != nil {
		return false
	}
	if err != nil {
		return isRetryableStorageError(err)
	}
	return isRetryableStorageStatus(resp.StatusCode)
}

func isRetryableStorageStatus(status int) bool {
	switch status {
	case http.StatusRequestTimeout,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func isRetryableStorageError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		if errors.Is(urlErr.Err, context.Canceled) || errors.Is(urlErr.Err, context.DeadlineExceeded) {
			return false
		}
		return true
	}
	return true
}

func storageRetryDelay(cfg RetryConfig, attempt int, resp *http.Response) time.Duration {
	if resp != nil {
		if wait, ok := parseRetryAfter(resp.Header.Get("Retry-After")); ok {
			return wait
		}
	}
	backoff := float64(cfg.InitialBackoff) * math.Pow(2, float64(attempt))
	delay := time.Duration(backoff)
	if delay > cfg.MaxBackoff {
		return cfg.MaxBackoff
	}
	return delay
}

func parseRetryAfter(value string) (time.Duration, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	if seconds, err := strconv.Atoi(value); err == nil && seconds >= 0 {
		return time.Duration(seconds) * time.Second, true
	}
	if when, err := http.ParseTime(value); err == nil {
		wait := time.Until(when)
		if wait < 0 {
			return 0, true
		}
		return wait, true
	}
	return 0, false
}

func sleepWithContext(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
