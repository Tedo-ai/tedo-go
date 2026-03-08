package tedo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		sep := "?"
		if params.Prefix != "" {
			path += sep + "prefix=" + params.Prefix
			sep = "&"
		}
		if params.Limit > 0 {
			path += sep + "limit=" + strconv.Itoa(params.Limit)
			sep = "&"
		}
		if params.Offset > 0 {
			path += sep + "offset=" + strconv.Itoa(params.Offset)
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
	path := fmt.Sprintf("/storage/v1/buckets/%s/objects/%s", bucketID, key)
	url := s.client.baseURL + path

	req, err := http.NewRequestWithContext(ctx, "PUT", url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.client.apiKey)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	} else {
		req.Header.Set("Content-Type", "application/octet-stream")
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

// GetObject downloads an object. The caller must close the returned ReadCloser.
func (s *StorageService) GetObject(ctx context.Context, bucketID, key string) (io.ReadCloser, string, error) {
	path := fmt.Sprintf("/storage/v1/buckets/%s/objects/%s", bucketID, key)
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
	return s.client.request(ctx, "DELETE", fmt.Sprintf("/storage/v1/buckets/%s/objects/%s", bucketID, key), nil, nil)
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
