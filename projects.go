package tedo

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ProjectsService handles Projects API calls.
type ProjectsService struct {
	client *Client
}

const (
	ProjectStatusCategoryStart      = "start"
	ProjectStatusCategoryInProgress = "in_progress"
	ProjectStatusCategoryCompleted  = "completed"
	ProjectStatusCategoryCanceled   = "canceled"
)

const (
	ProjectPriorityNone   = 0
	ProjectPriorityLow    = 1
	ProjectPriorityMedium = 2
	ProjectPriorityHigh   = 3
	ProjectPriorityUrgent = 4
)

// Project represents a Projects project.
type Project struct {
	ID          string    `json:"id"`
	TeamID      *string   `json:"team_id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Archived    bool      `json:"archived"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// WorkItem represents a Projects work item.
type WorkItem struct {
	ID             string     `json:"id"`
	DisplayID      string     `json:"display_id"`
	SequenceNumber int        `json:"sequence_number"`
	ProjectID      *string    `json:"project_id,omitempty"`
	TeamID         *string    `json:"team_id,omitempty"`
	ParentID       *string    `json:"parent_id,omitempty"`
	WorkItemTypeID *string    `json:"work_item_type_id,omitempty"`
	StatusID       *string    `json:"status_id,omitempty"`
	Title          string     `json:"title"`
	Description    string     `json:"description,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	ArchivedAt     *time.Time `json:"archived_at,omitempty"`
	AssigneeID     *string    `json:"assignee_id,omitempty"`
	Priority       *int       `json:"priority,omitempty"`
	DueDate        *string    `json:"due_date,omitempty"`
	Position       *float64   `json:"position,omitempty"`
	ChildCount     int        `json:"child_count"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// WorkflowStatus represents a work-item workflow status.
type WorkflowStatus struct {
	ID             string    `json:"id"`
	WorkItemTypeID *string   `json:"work_item_type_id,omitempty"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	Position       int       `json:"position"`
	IsDefault      bool      `json:"is_default"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// WorkItemType represents a configurable work-item type.
type WorkItemType struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	SingularName  *string   `json:"singular_name,omitempty"`
	PluralName    *string   `json:"plural_name,omitempty"`
	ParentTypeID  *string   `json:"parent_type_id,omitempty"`
	Prefix        *string   `json:"prefix,omitempty"`
	Color         *string   `json:"color,omitempty"`
	Position      int       `json:"position"`
	ShowInSidebar bool      `json:"show_in_sidebar"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PriorityLevel represents a configurable priority level.
type PriorityLevel struct {
	ID        string    `json:"id"`
	Level     int       `json:"level"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProjectComment represents a read-only Projects comment.
type ProjectComment struct {
	ID         string    `json:"id"`
	WorkItemID string    `json:"work_item_id"`
	ActorType  string    `json:"actor_type"`
	ActorRef   string    `json:"actor_ref"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// WorkItemActivity represents an activity feed entry.
type WorkItemActivity struct {
	ID        string         `json:"id"`
	ActorType string         `json:"actor_type"`
	ActorRef  string         `json:"actor_ref"`
	Action    string         `json:"action"`
	Data      map[string]any `json:"data,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

// ProjectAttachment represents a file reference attached to a work item.
type ProjectAttachment struct {
	ID           string    `json:"id"`
	WorkItemID   string    `json:"work_item_id"`
	FileID       string    `json:"file_id"`
	Position     int       `json:"position"`
	Filename     string    `json:"filename,omitempty"`
	MimeType     string    `json:"mime_type,omitempty"`
	Size         int64     `json:"size,omitempty"`
	Title        string    `json:"title,omitempty"`
	AltText      string    `json:"alt_text,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	ActorType    *string   `json:"actor_type,omitempty"`
	CreatedByRef *string   `json:"created_by_ref,omitempty"`
}

type ProjectList struct {
	Items      []Project `json:"items"`
	NextCursor string    `json:"next_cursor,omitempty"`
	HasMore    bool      `json:"has_more"`
}

type WorkItemList struct {
	Items      []WorkItem `json:"items"`
	NextCursor string     `json:"next_cursor,omitempty"`
	HasMore    bool       `json:"has_more"`
}

type WorkflowStatusList struct {
	Items      []WorkflowStatus `json:"items"`
	NextCursor string           `json:"next_cursor,omitempty"`
	HasMore    bool             `json:"has_more"`
}

type WorkItemTypeList struct {
	Items      []WorkItemType `json:"items"`
	NextCursor string         `json:"next_cursor,omitempty"`
	HasMore    bool           `json:"has_more"`
}

type PriorityLevelList struct {
	Items      []PriorityLevel `json:"items"`
	NextCursor string          `json:"next_cursor,omitempty"`
	HasMore    bool            `json:"has_more"`
}

type ProjectCommentList struct {
	Items      []ProjectComment `json:"items"`
	NextCursor string           `json:"next_cursor,omitempty"`
	HasMore    bool             `json:"has_more"`
}

type WorkItemActivityList struct {
	Items      []WorkItemActivity `json:"items"`
	NextCursor string             `json:"next_cursor,omitempty"`
	HasMore    bool               `json:"has_more"`
}

type ProjectAttachmentList struct {
	Items      []ProjectAttachment `json:"items"`
	NextCursor string              `json:"next_cursor,omitempty"`
	HasMore    bool                `json:"has_more"`
}

type DeleteResult struct {
	Deleted bool   `json:"deleted"`
	ID      string `json:"id"`
}

type ListProjectsParams struct {
	IncludeArchived bool
	Limit           int
	Cursor          string
}

type CreateProjectParams struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	TeamID      *string `json:"team_id,omitempty"`
}

type UpdateProjectParams struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	TeamID      *string `json:"team_id,omitempty"`
}

type ListWorkItemsParams struct {
	ProjectID        string
	WorkItemTypeID   string
	StatusID         string
	ParentID         string
	AssigneeID       string
	Priority         *int
	IncludeCompleted bool
	IncludeArchived  bool
	Limit            int
	Cursor           string
}

type ListPageParams struct {
	Limit  int
	Cursor string
}

type CreateWorkItemParams struct {
	ProjectID      *string  `json:"project_id,omitempty"`
	TeamID         *string  `json:"team_id,omitempty"`
	ParentID       *string  `json:"parent_id,omitempty"`
	WorkItemTypeID *string  `json:"work_item_type_id,omitempty"`
	StatusID       *string  `json:"status_id,omitempty"`
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	AssigneeID     *string  `json:"assignee_id,omitempty"`
	Priority       *int     `json:"priority,omitempty"`
	DueDate        *string  `json:"due_date,omitempty"`
	Position       *float64 `json:"position,omitempty"`
}

type UpdateWorkItemParams struct {
	ProjectID      *string  `json:"project_id,omitempty"`
	TeamID         *string  `json:"team_id,omitempty"`
	ParentID       *string  `json:"parent_id,omitempty"`
	WorkItemTypeID *string  `json:"work_item_type_id,omitempty"`
	StatusID       *string  `json:"status_id,omitempty"`
	Title          *string  `json:"title,omitempty"`
	Description    *string  `json:"description,omitempty"`
	AssigneeID     *string  `json:"assignee_id,omitempty"`
	Priority       *int     `json:"priority,omitempty"`
	DueDate        *string  `json:"due_date,omitempty"`
	Position       *float64 `json:"position,omitempty"`
}

type ListActivityParams struct {
	IncludeSubtasks bool
	IncludeComments bool
	Limit           int
	Cursor          string
}

type PeekNextDisplayIDParams struct {
	WorkItemTypeID string
}

type ListStatusesParams struct {
	WorkItemTypeID string
}

type CreateStatusParams struct {
	WorkItemTypeID *string `json:"work_item_type_id,omitempty"`
	Name           string  `json:"name"`
	Category       string  `json:"category"`
	Position       *int    `json:"position,omitempty"`
	IsDefault      *bool   `json:"is_default,omitempty"`
}

type UpdateStatusParams struct {
	WorkItemTypeID *string `json:"work_item_type_id,omitempty"`
	Name           *string `json:"name,omitempty"`
	Category       *string `json:"category,omitempty"`
	Position       *int    `json:"position,omitempty"`
	IsDefault      *bool   `json:"is_default,omitempty"`
}

type CreateWorkItemTypeParams struct {
	Name          string  `json:"name"`
	SingularName  string  `json:"singular_name,omitempty"`
	PluralName    string  `json:"plural_name,omitempty"`
	ParentTypeID  *string `json:"parent_type_id,omitempty"`
	Prefix        string  `json:"prefix,omitempty"`
	Color         string  `json:"color,omitempty"`
	Position      *int    `json:"position,omitempty"`
	ShowInSidebar *bool   `json:"show_in_sidebar,omitempty"`
}

type UpdateWorkItemTypeParams struct {
	Name          *string `json:"name,omitempty"`
	SingularName  *string `json:"singular_name,omitempty"`
	PluralName    *string `json:"plural_name,omitempty"`
	ParentTypeID  *string `json:"parent_type_id,omitempty"`
	Prefix        *string `json:"prefix,omitempty"`
	Color         *string `json:"color,omitempty"`
	Position      *int    `json:"position,omitempty"`
	ShowInSidebar *bool   `json:"show_in_sidebar,omitempty"`
}

type UpdatePriorityLevelParams struct {
	Name string `json:"name,omitempty"`
	Icon string `json:"icon,omitempty"`
}

type AttachFileParams struct {
	FileID      string `json:"file_id"`
	DisplayName string `json:"display_name,omitempty"`
}

type NextDisplayID struct {
	DisplayID string `json:"display_id"`
}

func (s *ProjectsService) ListProjects(ctx context.Context, params *ListProjectsParams) (*ProjectList, error) {
	var list ProjectList
	if err := s.client.request(ctx, http.MethodGet, "/projects/v1/projects"+listProjectsQuery(params), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (s *ProjectsService) CreateProject(ctx context.Context, params *CreateProjectParams, opts ...RequestOption) (*Project, error) {
	var project Project
	if err := s.client.requestWithOptions(ctx, http.MethodPost, "/projects/v1/projects", params, &project, requiredIdempotencyOptions(opts)...); err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *ProjectsService) GetProject(ctx context.Context, projectID string) (*Project, error) {
	var project Project
	if err := s.client.request(ctx, http.MethodGet, "/projects/v1/projects/"+pathEscape(projectID), nil, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *ProjectsService) UpdateProject(ctx context.Context, projectID string, params *UpdateProjectParams, opts ...RequestOption) (*Project, error) {
	var project Project
	if err := s.client.requestWithOptions(ctx, http.MethodPatch, "/projects/v1/projects/"+pathEscape(projectID), params, &project, opts...); err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *ProjectsService) ArchiveProject(ctx context.Context, projectID string, opts ...RequestOption) (*Project, error) {
	return s.projectAction(ctx, projectID, "archive", opts...)
}

func (s *ProjectsService) RestoreProject(ctx context.Context, projectID string, opts ...RequestOption) (*Project, error) {
	return s.projectAction(ctx, projectID, "restore", opts...)
}

func (s *ProjectsService) DeleteProject(ctx context.Context, projectID string, opts ...RequestOption) (*DeleteResult, error) {
	var result DeleteResult
	if err := s.client.requestWithOptions(ctx, http.MethodDelete, "/projects/v1/projects/"+pathEscape(projectID), nil, &result, requiredIdempotencyOptions(opts)...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ProjectsService) ListProjectWorkItems(ctx context.Context, projectID string, params *ListWorkItemsParams) (*WorkItemList, error) {
	var list WorkItemList
	path := "/projects/v1/projects/" + pathEscape(projectID) + "/work-items" + workItemListQuery(params, false)
	if err := s.client.request(ctx, http.MethodGet, path, nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (s *ProjectsService) CreateProjectWorkItem(ctx context.Context, projectID string, params *CreateWorkItemParams, opts ...RequestOption) (*WorkItem, error) {
	if params == nil {
		params = &CreateWorkItemParams{}
	}
	body := *params
	body.ProjectID = &projectID
	var item WorkItem
	if err := s.client.requestWithOptions(ctx, http.MethodPost, "/projects/v1/projects/"+pathEscape(projectID)+"/work-items", &body, &item, requiredIdempotencyOptions(opts)...); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ProjectsService) ListWorkItems(ctx context.Context, params *ListWorkItemsParams) (*WorkItemList, error) {
	var list WorkItemList
	if err := s.client.request(ctx, http.MethodGet, "/projects/v1/work-items"+workItemListQuery(params, true), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (s *ProjectsService) CreateWorkItem(ctx context.Context, params *CreateWorkItemParams, opts ...RequestOption) (*WorkItem, error) {
	var item WorkItem
	if err := s.client.requestWithOptions(ctx, http.MethodPost, "/projects/v1/work-items", params, &item, requiredIdempotencyOptions(opts)...); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ProjectsService) PeekNextDisplayID(ctx context.Context, params *PeekNextDisplayIDParams) (*NextDisplayID, error) {
	var result NextDisplayID
	query := url.Values{}
	if params != nil && params.WorkItemTypeID != "" {
		query.Set("work_item_type_id", params.WorkItemTypeID)
	}
	if err := s.client.request(ctx, http.MethodGet, "/projects/v1/work-items/next-display-id"+encodeQuery(query), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ProjectsService) GetWorkItem(ctx context.Context, workItemID string) (*WorkItem, error) {
	var item WorkItem
	if err := s.client.request(ctx, http.MethodGet, "/projects/v1/work-items/"+pathEscape(workItemID), nil, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ProjectsService) UpdateWorkItem(ctx context.Context, workItemID string, params *UpdateWorkItemParams, opts ...RequestOption) (*WorkItem, error) {
	var item WorkItem
	if err := s.client.requestWithOptions(ctx, http.MethodPatch, "/projects/v1/work-items/"+pathEscape(workItemID), params, &item, opts...); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ProjectsService) CompleteWorkItem(ctx context.Context, workItemID string, completed bool, opts ...RequestOption) (*WorkItem, error) {
	var item WorkItem
	body := map[string]bool{"completed": completed}
	if err := s.client.requestWithOptions(ctx, http.MethodPost, "/projects/v1/work-items/"+pathEscape(workItemID)+"/complete", body, &item, opts...); err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ProjectsService) ArchiveWorkItem(ctx context.Context, workItemID string, opts ...RequestOption) (*WorkItem, error) {
	return s.workItemAction(ctx, workItemID, "archive", opts...)
}

func (s *ProjectsService) RestoreWorkItem(ctx context.Context, workItemID string, opts ...RequestOption) (*WorkItem, error) {
	return s.workItemAction(ctx, workItemID, "restore", opts...)
}

func (s *ProjectsService) DeleteWorkItem(ctx context.Context, workItemID string, opts ...RequestOption) (*DeleteResult, error) {
	var result DeleteResult
	if err := s.client.requestWithOptions(ctx, http.MethodDelete, "/projects/v1/work-items/"+pathEscape(workItemID), nil, &result, requiredIdempotencyOptions(opts)...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ProjectsService) ListSubtasks(ctx context.Context, workItemID string, params *ListPageParams) (*WorkItemList, error) {
	var list WorkItemList
	if err := s.client.request(ctx, http.MethodGet, "/projects/v1/work-items/"+pathEscape(workItemID)+"/subtasks"+pageQuery(params), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (s *ProjectsService) ListWorkItemActivity(ctx context.Context, workItemID string, params *ListActivityParams) (*WorkItemActivityList, error) {
	var list WorkItemActivityList
	if err := s.client.request(ctx, http.MethodGet, "/projects/v1/work-items/"+pathEscape(workItemID)+"/activity"+activityQuery(params), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (s *ProjectsService) ListStatuses(ctx context.Context, params *ListStatusesParams) (*WorkflowStatusList, error) {
	var list WorkflowStatusList
	query := url.Values{}
	if params != nil && params.WorkItemTypeID != "" {
		query.Set("work_item_type_id", params.WorkItemTypeID)
	}
	if err := s.client.request(ctx, http.MethodGet, "/projects/v1/statuses"+encodeQuery(query), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (s *ProjectsService) CreateStatus(ctx context.Context, params *CreateStatusParams, opts ...RequestOption) (*WorkflowStatus, error) {
	var status WorkflowStatus
	if err := s.client.requestWithOptions(ctx, http.MethodPost, "/projects/v1/statuses", params, &status, requiredIdempotencyOptions(opts)...); err != nil {
		return nil, err
	}
	return &status, nil
}

func (s *ProjectsService) UpdateStatus(ctx context.Context, statusID string, params *UpdateStatusParams, opts ...RequestOption) (*WorkflowStatus, error) {
	var status WorkflowStatus
	if err := s.client.requestWithOptions(ctx, http.MethodPatch, "/projects/v1/statuses/"+pathEscape(statusID), params, &status, opts...); err != nil {
		return nil, err
	}
	return &status, nil
}

func (s *ProjectsService) DeleteStatus(ctx context.Context, statusID string, opts ...RequestOption) (*DeleteResult, error) {
	var result DeleteResult
	if err := s.client.requestWithOptions(ctx, http.MethodDelete, "/projects/v1/statuses/"+pathEscape(statusID), nil, &result, requiredIdempotencyOptions(opts)...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ProjectsService) ListWorkItemTypes(ctx context.Context) (*WorkItemTypeList, error) {
	var list WorkItemTypeList
	if err := s.client.request(ctx, http.MethodGet, "/projects/v1/work-item-types", nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (s *ProjectsService) CreateWorkItemType(ctx context.Context, params *CreateWorkItemTypeParams, opts ...RequestOption) (*WorkItemType, error) {
	var itemType WorkItemType
	if err := s.client.requestWithOptions(ctx, http.MethodPost, "/projects/v1/work-item-types", params, &itemType, requiredIdempotencyOptions(opts)...); err != nil {
		return nil, err
	}
	return &itemType, nil
}

func (s *ProjectsService) UpdateWorkItemType(ctx context.Context, workItemTypeID string, params *UpdateWorkItemTypeParams, opts ...RequestOption) (*WorkItemType, error) {
	var itemType WorkItemType
	if err := s.client.requestWithOptions(ctx, http.MethodPatch, "/projects/v1/work-item-types/"+pathEscape(workItemTypeID), params, &itemType, opts...); err != nil {
		return nil, err
	}
	return &itemType, nil
}

func (s *ProjectsService) DeleteWorkItemType(ctx context.Context, workItemTypeID string, opts ...RequestOption) (*DeleteResult, error) {
	var result DeleteResult
	if err := s.client.requestWithOptions(ctx, http.MethodDelete, "/projects/v1/work-item-types/"+pathEscape(workItemTypeID), nil, &result, requiredIdempotencyOptions(opts)...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ProjectsService) ListPriorityLevels(ctx context.Context) (*PriorityLevelList, error) {
	var list PriorityLevelList
	if err := s.client.request(ctx, http.MethodGet, "/projects/v1/priority-levels", nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (s *ProjectsService) UpdatePriorityLevel(ctx context.Context, level int, params *UpdatePriorityLevelParams, opts ...RequestOption) (*PriorityLevel, error) {
	var priority PriorityLevel
	if err := s.client.requestWithOptions(ctx, http.MethodPatch, "/projects/v1/priority-levels/"+strconv.Itoa(level), params, &priority, opts...); err != nil {
		return nil, err
	}
	return &priority, nil
}

func (s *ProjectsService) ResetPriorityLevel(ctx context.Context, level int, opts ...RequestOption) (*PriorityLevel, error) {
	var priority PriorityLevel
	if err := s.client.requestWithOptions(ctx, http.MethodPost, "/projects/v1/priority-levels/"+strconv.Itoa(level)+"/reset", nil, &priority, opts...); err != nil {
		return nil, err
	}
	return &priority, nil
}

func (s *ProjectsService) ListComments(ctx context.Context, workItemID string, params *ListPageParams) (*ProjectCommentList, error) {
	var list ProjectCommentList
	if err := s.client.request(ctx, http.MethodGet, "/projects/v1/work-items/"+pathEscape(workItemID)+"/comments"+pageQuery(params), nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (s *ProjectsService) ListAttachments(ctx context.Context, workItemID string) (*ProjectAttachmentList, error) {
	var list ProjectAttachmentList
	if err := s.client.request(ctx, http.MethodGet, "/projects/v1/work-items/"+pathEscape(workItemID)+"/attachments", nil, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (s *ProjectsService) AttachFile(ctx context.Context, workItemID string, params *AttachFileParams, opts ...RequestOption) (*ProjectAttachment, error) {
	var attachment ProjectAttachment
	if err := s.client.requestWithOptions(ctx, http.MethodPost, "/projects/v1/work-items/"+pathEscape(workItemID)+"/attachments", params, &attachment, requiredIdempotencyOptions(opts)...); err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (s *ProjectsService) DetachAttachment(ctx context.Context, workItemID, attachmentID string, opts ...RequestOption) (*DeleteResult, error) {
	var result DeleteResult
	path := "/projects/v1/work-items/" + pathEscape(workItemID) + "/attachments/" + pathEscape(attachmentID)
	if err := s.client.requestWithOptions(ctx, http.MethodDelete, path, nil, &result, requiredIdempotencyOptions(opts)...); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ProjectsService) projectAction(ctx context.Context, projectID, action string, opts ...RequestOption) (*Project, error) {
	var project Project
	path := "/projects/v1/projects/" + pathEscape(projectID) + "/" + action
	if err := s.client.requestWithOptions(ctx, http.MethodPost, path, nil, &project, opts...); err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *ProjectsService) workItemAction(ctx context.Context, workItemID, action string, opts ...RequestOption) (*WorkItem, error) {
	var item WorkItem
	path := "/projects/v1/work-items/" + pathEscape(workItemID) + "/" + action
	if err := s.client.requestWithOptions(ctx, http.MethodPost, path, nil, &item, opts...); err != nil {
		return nil, err
	}
	return &item, nil
}

func listProjectsQuery(params *ListProjectsParams) string {
	query := url.Values{}
	if params == nil {
		return ""
	}
	if params.IncludeArchived {
		query.Set("include_archived", "true")
	}
	addPageQuery(query, params.Limit, params.Cursor)
	return encodeQuery(query)
}

func workItemListQuery(params *ListWorkItemsParams, includeProjectID bool) string {
	query := url.Values{}
	if params == nil {
		return ""
	}
	if includeProjectID && params.ProjectID != "" {
		query.Set("project_id", params.ProjectID)
	}
	if params.WorkItemTypeID != "" {
		query.Set("work_item_type_id", params.WorkItemTypeID)
	}
	if params.StatusID != "" {
		query.Set("status_id", params.StatusID)
	}
	if params.ParentID != "" {
		query.Set("parent_id", params.ParentID)
	}
	if params.AssigneeID != "" {
		query.Set("assignee_id", params.AssigneeID)
	}
	if params.Priority != nil {
		query.Set("priority", strconv.Itoa(*params.Priority))
	}
	if params.IncludeCompleted {
		query.Set("include_completed", "true")
	}
	if params.IncludeArchived {
		query.Set("include_archived", "true")
	}
	addPageQuery(query, params.Limit, params.Cursor)
	return encodeQuery(query)
}

func pageQuery(params *ListPageParams) string {
	query := url.Values{}
	if params != nil {
		addPageQuery(query, params.Limit, params.Cursor)
	}
	return encodeQuery(query)
}

func activityQuery(params *ListActivityParams) string {
	query := url.Values{}
	if params != nil {
		if params.IncludeSubtasks {
			query.Set("include_subtasks", "true")
		}
		if params.IncludeComments {
			query.Set("include_comments", "true")
		}
		addPageQuery(query, params.Limit, params.Cursor)
	}
	return encodeQuery(query)
}

func addPageQuery(query url.Values, limit int, cursor string) {
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		query.Set("cursor", cursor)
	}
}

func encodeQuery(query url.Values) string {
	if len(query) == 0 {
		return ""
	}
	return "?" + query.Encode()
}

func pathEscape(value string) string {
	return url.PathEscape(value)
}

func requiredIdempotencyOptions(opts []RequestOption) []RequestOption {
	out := make([]RequestOption, 0, len(opts)+1)
	out = append(out, WithIdempotencyKey(newIdempotencyKey()))
	out = append(out, opts...)
	return out
}

func newIdempotencyKey() string {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err == nil {
		return "tedo_go_" + hex.EncodeToString(buf[:])
	}
	return fmt.Sprintf("tedo_go_%d", time.Now().UnixNano())
}
