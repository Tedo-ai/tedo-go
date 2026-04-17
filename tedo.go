// Package tedo provides a Go client for the Tedo API.
//
// Usage:
//
//	client := tedo.NewClient("tedo_live_xxx")
//	customer, err := client.Billing.CreateCustomer(ctx, &tedo.CreateCustomerParams{
//	    Email: "user@example.com",
//	    Name:  "Acme Corp",
//	})
package tedo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.tedo.ai"
)

var defaultRetryConfig = RetryConfig{
	MaxRetries:     3,
	InitialBackoff: 200 * time.Millisecond,
	MaxBackoff:     2 * time.Second,
}

// RetryConfig configures how retryable requests are retried.
type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

// Client is the Tedo API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	retry      RetryConfig

	// Services
	Billing *BillingService
	Sales   *SalesService
	Storage *StorageService
}

// NewClient creates a new Tedo API client.
func NewClient(apiKey string) *Client {
	c := &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{},
		retry:      defaultRetryConfig,
	}

	// Initialize services
	c.Billing = &BillingService{client: c}
	c.Sales = &SalesService{client: c}
	c.Storage = &StorageService{client: c}

	return c
}

// WithBaseURL sets a custom base URL (useful for testing).
func (c *Client) WithBaseURL(url string) *Client {
	c.baseURL = url
	return c
}

// WithHTTPClient sets a custom HTTP client.
func (c *Client) WithHTTPClient(httpClient *http.Client) *Client {
	c.httpClient = httpClient
	return c
}

// WithRetryConfig sets the retry policy used by retry-aware SDK operations.
func (c *Client) WithRetryConfig(cfg RetryConfig) *Client {
	c.retry = normalizeRetryConfig(cfg)
	return c
}

func normalizeRetryConfig(cfg RetryConfig) RetryConfig {
	if cfg.MaxRetries < 0 {
		cfg.MaxRetries = 0
	}
	if cfg.InitialBackoff <= 0 {
		cfg.InitialBackoff = defaultRetryConfig.InitialBackoff
	}
	if cfg.MaxBackoff <= 0 {
		cfg.MaxBackoff = defaultRetryConfig.MaxBackoff
	}
	if cfg.MaxBackoff < cfg.InitialBackoff {
		cfg.MaxBackoff = cfg.InitialBackoff
	}
	return cfg
}

// request performs an API request and decodes the response.
func (c *Client) request(ctx context.Context, method, path string, body, result any) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		return parseError(resp.StatusCode, respBody)
	}

	// Decode successful response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// Error types

// Error represents an API error.
type Error struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Field      string `json:"field,omitempty"`
}

func (e *Error) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("tedo: %s - %s (field: %s)", e.Code, e.Message, e.Field)
	}
	return fmt.Sprintf("tedo: %s - %s", e.Code, e.Message)
}

// IsNotFound returns true if the error is a 404 Not Found.
func IsNotFound(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == 404
	}
	return false
}

// IsValidationError returns true if the error is a 400 Bad Request.
func IsValidationError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == 400
	}
	return false
}

// IsUnauthorized returns true if the error is a 401 Unauthorized.
func IsUnauthorized(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.StatusCode == 401
	}
	return false
}

func parseError(statusCode int, body []byte) error {
	var apiErr Error
	if err := json.Unmarshal(body, &apiErr); err != nil {
		// If we can't parse the error, create a generic one
		apiErr = Error{
			Code:    "unknown_error",
			Message: string(body),
		}
	}
	apiErr.StatusCode = statusCode
	return &apiErr
}
