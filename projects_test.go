package tedo

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestNewClientInitializesProjectsService(t *testing.T) {
	client := NewClient("tedo_live_test")
	if client.Projects == nil {
		t.Fatal("Projects service is nil")
	}
}

func TestProjectsService_CreateProjectSendsIdempotencyKey(t *testing.T) {
	client := newProjectsTestClient(t, func(r *http.Request) *http.Response {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if got, want := r.URL.Path, "/projects/v1/projects"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("Authorization"), "Bearer tedo_live_test"; got != want {
			t.Fatalf("authorization = %q, want %q", got, want)
		}
		if key := r.Header.Get("Idempotency-Key"); key == "" {
			t.Fatal("missing Idempotency-Key")
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if _, ok := body["idempotency_key"]; ok {
			t.Fatal("idempotency key leaked into JSON body")
		}
		if got, want := body["name"], "Launch"; got != want {
			t.Fatalf("name = %v, want %v", got, want)
		}

		return jsonResponse(http.StatusCreated, map[string]any{
			"id":          "proj-1",
			"name":        "Launch",
			"description": "Q2",
			"archived":    false,
			"created_at":  "2026-05-08T00:00:00Z",
			"updated_at":  "2026-05-08T00:00:00Z",
		})
	})

	project, err := client.Projects.CreateProject(context.Background(), &CreateProjectParams{Name: "Launch", Description: "Q2"})
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}
	if project.ID != "proj-1" {
		t.Fatalf("project ID = %q, want proj-1", project.ID)
	}
}

func TestProjectsService_RequestOptionsOverrideIdempotencyKeyAndRequestID(t *testing.T) {
	client := newProjectsTestClient(t, func(r *http.Request) *http.Response {
		if got, want := r.Header.Get("Idempotency-Key"), "idem_123"; got != want {
			t.Fatalf("Idempotency-Key = %q, want %q", got, want)
		}
		if got, want := r.Header.Get("X-Request-ID"), "req_123"; got != want {
			t.Fatalf("X-Request-ID = %q, want %q", got, want)
		}
		return jsonResponse(http.StatusOK, map[string]any{
			"deleted": true,
			"id":      "proj-1",
		})
	})

	result, err := client.Projects.DeleteProject(context.Background(), "proj-1", WithIdempotencyKey("idem_123"), WithRequestID("req_123"))
	if err != nil {
		t.Fatalf("DeleteProject failed: %v", err)
	}
	if !result.Deleted || result.ID != "proj-1" {
		t.Fatalf("delete result = %#v", result)
	}
}

func TestProjectsService_ListWorkItemsEncodesCursorFilters(t *testing.T) {
	client := newProjectsTestClient(t, func(r *http.Request) *http.Response {
		if got, want := r.URL.Path, "/projects/v1/work-items"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		wantQuery := "cursor=eyJvZmZzZXQiOjUwfQ&include_archived=true&include_completed=true&limit=50&priority=2&project_id=proj-1&status_id=status-1"
		if got := r.URL.RawQuery; got != wantQuery {
			t.Fatalf("query = %q, want %q", got, wantQuery)
		}
		return jsonResponse(http.StatusOK, map[string]any{
			"items":       []map[string]any{},
			"next_cursor": nil,
			"has_more":    false,
		})
	})

	priority := ProjectPriorityMedium
	list, err := client.Projects.ListWorkItems(context.Background(), &ListWorkItemsParams{
		ProjectID:        "proj-1",
		StatusID:         "status-1",
		Priority:         &priority,
		IncludeCompleted: true,
		IncludeArchived:  true,
		Limit:            50,
		Cursor:           "eyJvZmZzZXQiOjUwfQ",
	})
	if err != nil {
		t.Fatalf("ListWorkItems failed: %v", err)
	}
	if list.HasMore {
		t.Fatal("HasMore = true, want false")
	}
}

func TestProjectsService_ListCommentsUsesReadOnlyEndpoint(t *testing.T) {
	client := newProjectsTestClient(t, func(r *http.Request) *http.Response {
		if got, want := r.Method, http.MethodGet; got != want {
			t.Fatalf("method = %s, want %s", got, want)
		}
		if strings.Contains(r.URL.Path, "/comments/") {
			t.Fatalf("unexpected comment write-style path: %s", r.URL.Path)
		}
		if got, want := r.URL.Path, "/projects/v1/work-items/work-1/comments"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		return jsonResponse(http.StatusOK, map[string]any{
			"items": []map[string]any{
				{
					"id":           "comment-1",
					"work_item_id": "work-1",
					"actor_type":   "api_key",
					"actor_ref":    "api_key:key-1",
					"content":      "ready",
					"created_at":   "2026-05-08T00:00:00Z",
					"updated_at":   "2026-05-08T00:00:00Z",
				},
			},
			"has_more": false,
		})
	})

	list, err := client.Projects.ListComments(context.Background(), "work-1", nil)
	if err != nil {
		t.Fatalf("ListComments failed: %v", err)
	}
	if len(list.Items) != 1 || list.Items[0].ActorRef != "api_key:key-1" {
		t.Fatalf("comments = %#v", list.Items)
	}
}

func TestParseErrorCapturesCanonicalProjectsEnvelope(t *testing.T) {
	err := parseError(http.StatusForbidden, []byte(`{
		"code":"permission_denied",
		"message":"API key lacks projects.projects.write",
		"details":{"permission":"projects.projects.write"},
		"request_id":"req_123"
	}`))

	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *Error", err)
	}
	if apiErr.Code != "permission_denied" || apiErr.RequestID != "req_123" {
		t.Fatalf("api error = %#v", apiErr)
	}
	if apiErr.Details["permission"] != "projects.projects.write" {
		t.Fatalf("details = %#v", apiErr.Details)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func newProjectsTestClient(t *testing.T, handler func(*http.Request) *http.Response) *Client {
	t.Helper()
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return handler(req), nil
	})}
	return NewClient("tedo_live_test").WithBaseURL("https://api.test").WithHTTPClient(httpClient)
}

func jsonResponse(status int, body any) *http.Response {
	var buf strings.Builder
	_ = json.NewEncoder(&buf).Encode(body)
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(buf.String())),
	}
}
