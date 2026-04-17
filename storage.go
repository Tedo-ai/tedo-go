package tedo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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
	if err := s.client.request(ctx, "GET", "/storage/v1/buckets", nil, &buckets); err != nil {
		return nil, err
	}
	return buckets, nil
}

// CreateBucket creates a new bucket.
func (s *StorageService) CreateBucket(ctx context.Context, params *CreateBucketParams) (*Bucket, error) {
	var bucket Bucket
	if err := s.client.request(ctx, "POST", "/storage/v1/buckets", params, &bucket); err != nil {
		return nil, err
	}
	return &bucket, nil
}

// GetBucket gets a bucket by ID.
func (s *StorageService) GetBucket(ctx context.Context, bucketID string) (*Bucket, error) {
	var bucket Bucket
	if err := s.client.request(ctx, "GET", fmt.Sprintf("/storage/v1/buckets/%s", bucketID), nil, &bucket); err != nil {
		return nil, err
	}
	return &bucket, nil
}

// DeleteBucket deletes an empty bucket.
func (s *StorageService) DeleteBucket(ctx context.Context, bucketID string) error {
	return s.client.request(ctx, "DELETE", fmt.Sprintf("/storage/v1/buckets/%s", bucketID), nil, nil)
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
	if err := s.client.request(ctx, "GET", path, nil, &list); err != nil {
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
	path := fmt.Sprintf("/storage/v1/buckets/%s/objects/%s", bucketID, url.PathEscape(key))
	url := s.client.baseURL + path

	req, err := http.NewRequestWithContext(ctx, "PUT", url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.client.apiKey)
	contentType := "application/octet-stream"
	if opts != nil && opts.ContentType != "" {
		contentType = opts.ContentType
	}
	req.Header.Set("Content-Type", contentType)
	if opts != nil && opts.ContentSHA256 != "" {
		req.Header.Set("X-Content-Sha256", opts.ContentSHA256)
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, parseError(resp.StatusCode, respBody)
	}

	var obj Object
	if err := json.Unmarshal(respBody, &obj); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &obj, nil
}

// HeadObject fetches object metadata without downloading the body.
func (s *StorageService) HeadObject(ctx context.Context, bucketID, key string) (*Object, error) {
	path := fmt.Sprintf("/storage/v1/buckets/%s/objects/%s", bucketID, url.PathEscape(key))
	url := s.client.baseURL + path

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.client.apiKey)

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, parseError(resp.StatusCode, []byte(resp.Status))
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
	url := s.client.baseURL + path

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.client.apiKey)

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("do request: %w", err)
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
	return s.client.request(ctx, "DELETE", fmt.Sprintf("/storage/v1/buckets/%s/objects/%s", bucketID, url.PathEscape(key)), nil, nil)
}

// --- Pre-signed URLs ---

// PresignURL creates a pre-signed URL for temporary public access.
func (s *StorageService) PresignURL(ctx context.Context, bucketID string, params *PresignParams) (*PresignResponse, error) {
	var resp PresignResponse
	if err := s.client.request(ctx, "POST", fmt.Sprintf("/storage/v1/buckets/%s/presign", bucketID), params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Usage ---

// GetUsage returns storage usage statistics.
func (s *StorageService) GetUsage(ctx context.Context) (*StorageUsage, error) {
	var usage StorageUsage
	if err := s.client.request(ctx, "GET", "/storage/v1/usage", nil, &usage); err != nil {
		return nil, err
	}
	return &usage, nil
}
