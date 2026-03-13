package tedo

import (
	"context"
	"fmt"
	"time"
)

// Activity types
const (
	ActivityTypeTask     = "task"
	ActivityTypeCall     = "call"
	ActivityTypeEmail    = "email"
	ActivityTypeMeeting  = "meeting"
	ActivityTypeDeadline = "deadline"
)

// Pipeline resource types
const (
	ResourceTypeLead = "lead"
	ResourceTypeDeal = "deal"
)

// Stage outcomes
const (
	OutcomePositive = "positive"
	OutcomeNegative = "negative"
)

// ActivityLink represents a link between an activity and a sales entity.
type ActivityLink struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	IsPrimary  bool   `json:"is_primary,omitempty"`
}

// LeadLink creates a link to a lead (primary by default).
func LeadLink(id string) ActivityLink {
	return ActivityLink{EntityType: "lead", EntityID: id, IsPrimary: true}
}

// DealLink creates a link to a deal (primary by default).
func DealLink(id string) ActivityLink {
	return ActivityLink{EntityType: "deal", EntityID: id, IsPrimary: true}
}

// PersonLink creates a link to a person.
func PersonLink(id string) ActivityLink {
	return ActivityLink{EntityType: "person", EntityID: id}
}

// OrganizationLink creates a link to an organization.
func OrganizationLink(id string) ActivityLink {
	return ActivityLink{EntityType: "organization", EntityID: id}
}

// LeadNoteLink creates a note link to a lead.
func LeadNoteLink(id string) SalesNoteLink {
	return SalesNoteLink{EntityType: "lead", EntityID: id, IsPrimary: true}
}

// DealNoteLink creates a note link to a deal.
func DealNoteLink(id string) SalesNoteLink {
	return SalesNoteLink{EntityType: "deal", EntityID: id, IsPrimary: true}
}

// PersonNoteLink creates a note link to a person.
func PersonNoteLink(id string) SalesNoteLink {
	return SalesNoteLink{EntityType: "person", EntityID: id}
}

// OrganizationNoteLink creates a note link to an organization.
func OrganizationNoteLink(id string) SalesNoteLink {
	return SalesNoteLink{EntityType: "organization", EntityID: id}
}

// SalesService handles sales-related API calls.
type SalesService struct {
	client *Client
}

// ============================================================
// PIPELINES
// ============================================================

// Pipeline represents a sales pipeline for leads or deals.
type Pipeline struct {
	ID           string    `json:"id"`
	WorkspaceID  string    `json:"workspace_id,omitempty"`
	Name         string    `json:"name"`
	ResourceType string    `json:"resource_type"` // "lead" or "deal"
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreatePipelineParams are the parameters for creating a pipeline.
type CreatePipelineParams struct {
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"` // "lead" or "deal"
}

// CreatePipeline creates a new sales pipeline.
func (s *SalesService) CreatePipeline(ctx context.Context, params *CreatePipelineParams) (*Pipeline, error) {
	var pipeline Pipeline
	err := s.client.request(ctx, "POST", "/sales/v1/pipelines", params, &pipeline)
	if err != nil {
		return nil, err
	}
	return &pipeline, nil
}

// PipelineList is a list of pipelines.
type PipelineList struct {
	Pipelines []Pipeline `json:"pipelines"`
	Total     int        `json:"total"`
}

// ListPipelines lists all pipelines.
func (s *SalesService) ListPipelines(ctx context.Context) (*PipelineList, error) {
	var list PipelineList
	err := s.client.request(ctx, "GET", "/sales/v1/pipelines", nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// GetPipeline retrieves a pipeline by ID.
func (s *SalesService) GetPipeline(ctx context.Context, id string) (*Pipeline, error) {
	var pipeline Pipeline
	err := s.client.request(ctx, "GET", "/sales/v1/pipelines/"+id, nil, &pipeline)
	if err != nil {
		return nil, err
	}
	return &pipeline, nil
}

// UpdatePipelineParams are the parameters for updating a pipeline.
type UpdatePipelineParams struct {
	Name         *string `json:"name,omitempty"`
	ResourceType *string `json:"resource_type,omitempty"`
}

// UpdatePipeline updates a pipeline.
func (s *SalesService) UpdatePipeline(ctx context.Context, id string, params *UpdatePipelineParams) (*Pipeline, error) {
	var pipeline Pipeline
	err := s.client.request(ctx, "PATCH", "/sales/v1/pipelines/"+id, params, &pipeline)
	if err != nil {
		return nil, err
	}
	return &pipeline, nil
}

// DeletePipeline deletes a pipeline.
func (s *SalesService) DeletePipeline(ctx context.Context, id string) error {
	return s.client.request(ctx, "DELETE", "/sales/v1/pipelines/"+id, nil, nil)
}

// ============================================================
// STAGES
// ============================================================

// PipelineStage represents a stage within a pipeline.
type PipelineStage struct {
	ID         string    `json:"id"`
	PipelineID string    `json:"pipeline_id"`
	Name       string    `json:"name"`
	Position   int       `json:"position"`
	IsTerminal bool      `json:"is_terminal"`
	Outcome    *string   `json:"outcome,omitempty"` // "positive", "negative", or nil
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateStageParams are the parameters for creating a pipeline stage.
type CreateStageParams struct {
	Name       string  `json:"name"`
	Position   int     `json:"position"`
	IsTerminal bool    `json:"is_terminal,omitempty"`
	Outcome    *string `json:"outcome,omitempty"`
}

// CreateStage creates a new stage in a pipeline.
func (s *SalesService) CreateStage(ctx context.Context, pipelineID string, params *CreateStageParams) (*PipelineStage, error) {
	var stage PipelineStage
	body := struct {
		PipelineID string  `json:"pipeline_id"`
		Name       string  `json:"name"`
		Position   int     `json:"position"`
		IsTerminal bool    `json:"is_terminal,omitempty"`
		Outcome    *string `json:"outcome,omitempty"`
	}{
		PipelineID: pipelineID,
		Name:       params.Name,
		Position:   params.Position,
		IsTerminal: params.IsTerminal,
		Outcome:    params.Outcome,
	}
	err := s.client.request(ctx, "POST", "/sales/v1/stages", body, &stage)
	if err != nil {
		return nil, err
	}
	return &stage, nil
}

// StageList is a list of pipeline stages.
type StageList struct {
	Stages []PipelineStage `json:"stages"`
	Total  int             `json:"total"`
}

// ListStages lists all stages for a pipeline.
func (s *SalesService) ListStages(ctx context.Context, pipelineID string) (*StageList, error) {
	var list StageList
	err := s.client.request(ctx, "GET", "/sales/v1/stages?pipeline_id="+pipelineID, nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// GetStage retrieves a stage by ID.
func (s *SalesService) GetStage(ctx context.Context, id string) (*PipelineStage, error) {
	var stage PipelineStage
	err := s.client.request(ctx, "GET", "/sales/v1/stages/"+id, nil, &stage)
	if err != nil {
		return nil, err
	}
	return &stage, nil
}

// UpdateStageParams are the parameters for updating a pipeline stage.
type UpdateStageParams struct {
	Name       *string `json:"name,omitempty"`
	Position   *int    `json:"position,omitempty"`
	IsTerminal *bool   `json:"is_terminal,omitempty"`
	Outcome    *string `json:"outcome,omitempty"`
}

// UpdateStage updates a pipeline stage.
func (s *SalesService) UpdateStage(ctx context.Context, id string, params *UpdateStageParams) (*PipelineStage, error) {
	var stage PipelineStage
	err := s.client.request(ctx, "PATCH", "/sales/v1/stages/"+id, params, &stage)
	if err != nil {
		return nil, err
	}
	return &stage, nil
}

// DeleteStage deletes a pipeline stage.
func (s *SalesService) DeleteStage(ctx context.Context, id string) error {
	return s.client.request(ctx, "DELETE", "/sales/v1/stages/"+id, nil, nil)
}

// ============================================================
// LEADS
// ============================================================

// Lead represents a sales lead in a pipeline.
type Lead struct {
	ID             string    `json:"id"`
	Label          string    `json:"label"`
	PipelineID     string    `json:"pipeline_id"`
	StageID        string    `json:"stage_id"`
	PersonID       *string   `json:"person_id,omitempty"`
	OrganizationID *string   `json:"organization_id,omitempty"`
	Source         *string   `json:"source,omitempty"`
	OwnerID        *string   `json:"owner_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CreateLeadParams are the parameters for creating a lead.
type CreateLeadParams struct {
	Label          string  `json:"label"`
	PipelineID     string  `json:"pipeline_id"`
	StageID        string  `json:"stage_id,omitempty"`
	PersonID       *string `json:"person_id,omitempty"`
	OrganizationID *string `json:"organization_id,omitempty"`
	Source         *string `json:"source,omitempty"`
	OwnerID        *string `json:"owner_id,omitempty"`
}

// CreateLead creates a new lead.
func (s *SalesService) CreateLead(ctx context.Context, params *CreateLeadParams) (*Lead, error) {
	var lead Lead
	err := s.client.request(ctx, "POST", "/sales/v1/leads", params, &lead)
	if err != nil {
		return nil, err
	}
	return &lead, nil
}

// ListLeadsParams are the parameters for listing leads.
type ListLeadsParams struct {
	PipelineID string `json:"pipeline_id,omitempty"`
}

// LeadList is a list of leads.
type LeadList struct {
	Leads []Lead `json:"leads"`
	Total int    `json:"total"`
}

// ListLeads lists leads with optional pipeline filter.
func (s *SalesService) ListLeads(ctx context.Context, params *ListLeadsParams) (*LeadList, error) {
	path := "/sales/v1/leads"
	if params != nil && params.PipelineID != "" {
		path += "?pipeline_id=" + params.PipelineID
	}
	var list LeadList
	err := s.client.request(ctx, "GET", path, nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// GetLead retrieves a lead by ID.
func (s *SalesService) GetLead(ctx context.Context, id string) (*Lead, error) {
	var lead Lead
	err := s.client.request(ctx, "GET", "/sales/v1/leads/"+id, nil, &lead)
	if err != nil {
		return nil, err
	}
	return &lead, nil
}

// UpdateLeadParams are the parameters for updating a lead.
type UpdateLeadParams struct {
	Label          *string `json:"label,omitempty"`
	PipelineID     *string `json:"pipeline_id,omitempty"`
	StageID        *string `json:"stage_id,omitempty"`
	PersonID       *string `json:"person_id,omitempty"`
	OrganizationID *string `json:"organization_id,omitempty"`
	Source         *string `json:"source,omitempty"`
	OwnerID        *string `json:"owner_id,omitempty"`
}

// UpdateLead updates a lead.
func (s *SalesService) UpdateLead(ctx context.Context, id string, params *UpdateLeadParams) (*Lead, error) {
	var lead Lead
	err := s.client.request(ctx, "PATCH", "/sales/v1/leads/"+id, params, &lead)
	if err != nil {
		return nil, err
	}
	return &lead, nil
}

// DeleteLead deletes a lead.
func (s *SalesService) DeleteLead(ctx context.Context, id string) error {
	return s.client.request(ctx, "DELETE", "/sales/v1/leads/"+id, nil, nil)
}

// MoveLeadStage moves a lead to a different pipeline stage.
func (s *SalesService) MoveLeadStage(ctx context.Context, id string, stageID string) (*Lead, error) {
	var lead Lead
	body := map[string]string{"stage_id": stageID}
	err := s.client.request(ctx, "POST", "/sales/v1/leads/"+id+"/move", body, &lead)
	if err != nil {
		return nil, err
	}
	return &lead, nil
}

// ConvertLeadParams are the parameters for converting a lead to a deal.
type ConvertLeadParams struct {
	DealPipelineID string   `json:"deal_pipeline_id"`
	DealStageID    string   `json:"deal_stage_id"`
	DealLabel      string   `json:"deal_label,omitempty"`
	Value          *float64 `json:"value,omitempty"`
	Currency       string   `json:"currency,omitempty"`
}

// ConvertLeadToDeal converts a lead into a deal.
func (s *SalesService) ConvertLeadToDeal(ctx context.Context, id string, params *ConvertLeadParams) (*Deal, error) {
	var deal Deal
	err := s.client.request(ctx, "POST", "/sales/v1/leads/"+id+"/convert", params, &deal)
	if err != nil {
		return nil, err
	}
	return &deal, nil
}

// ============================================================
// DEALS
// ============================================================

// Deal represents a sales deal in a pipeline.
type Deal struct {
	ID                string    `json:"id"`
	Label             string    `json:"label"`
	PipelineID        string    `json:"pipeline_id"`
	StageID           string    `json:"stage_id"`
	PersonID          *string   `json:"person_id,omitempty"`
	OrganizationID    *string   `json:"organization_id,omitempty"`
	Value             *float64  `json:"value,omitempty"`
	Currency          string    `json:"currency"`
	ExpectedCloseDate *string   `json:"expected_close_date,omitempty"`
	OwnerID           *string   `json:"owner_id,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CreateDealParams are the parameters for creating a deal.
type CreateDealParams struct {
	Label             string   `json:"label"`
	PipelineID        string   `json:"pipeline_id"`
	StageID           string   `json:"stage_id,omitempty"`
	PersonID          *string  `json:"person_id,omitempty"`
	OrganizationID    *string  `json:"organization_id,omitempty"`
	Value             *float64 `json:"value,omitempty"`
	Currency          string   `json:"currency,omitempty"`
	ExpectedCloseDate *string  `json:"expected_close_date,omitempty"`
	OwnerID           *string  `json:"owner_id,omitempty"`
}

// CreateDeal creates a new deal.
func (s *SalesService) CreateDeal(ctx context.Context, params *CreateDealParams) (*Deal, error) {
	var deal Deal
	err := s.client.request(ctx, "POST", "/sales/v1/deals", params, &deal)
	if err != nil {
		return nil, err
	}
	return &deal, nil
}

// ListDealsParams are the parameters for listing deals.
type ListDealsParams struct {
	PipelineID string `json:"pipeline_id,omitempty"`
}

// DealList is a list of deals.
type DealList struct {
	Deals []Deal `json:"deals"`
	Total int    `json:"total"`
}

// ListDeals lists deals with optional pipeline filter.
func (s *SalesService) ListDeals(ctx context.Context, params *ListDealsParams) (*DealList, error) {
	path := "/sales/v1/deals"
	if params != nil && params.PipelineID != "" {
		path += "?pipeline_id=" + params.PipelineID
	}
	var list DealList
	err := s.client.request(ctx, "GET", path, nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// GetDeal retrieves a deal by ID.
func (s *SalesService) GetDeal(ctx context.Context, id string) (*Deal, error) {
	var deal Deal
	err := s.client.request(ctx, "GET", "/sales/v1/deals/"+id, nil, &deal)
	if err != nil {
		return nil, err
	}
	return &deal, nil
}

// UpdateDealParams are the parameters for updating a deal.
type UpdateDealParams struct {
	Label             *string  `json:"label,omitempty"`
	PipelineID        *string  `json:"pipeline_id,omitempty"`
	StageID           *string  `json:"stage_id,omitempty"`
	PersonID          *string  `json:"person_id,omitempty"`
	OrganizationID    *string  `json:"organization_id,omitempty"`
	Value             *float64 `json:"value,omitempty"`
	Currency          *string  `json:"currency,omitempty"`
	ExpectedCloseDate *string  `json:"expected_close_date,omitempty"`
	OwnerID           *string  `json:"owner_id,omitempty"`
}

// UpdateDeal updates a deal.
func (s *SalesService) UpdateDeal(ctx context.Context, id string, params *UpdateDealParams) (*Deal, error) {
	var deal Deal
	err := s.client.request(ctx, "PATCH", "/sales/v1/deals/"+id, params, &deal)
	if err != nil {
		return nil, err
	}
	return &deal, nil
}

// DeleteDeal deletes a deal.
func (s *SalesService) DeleteDeal(ctx context.Context, id string) error {
	return s.client.request(ctx, "DELETE", "/sales/v1/deals/"+id, nil, nil)
}

// MoveDealStage moves a deal to a different pipeline stage.
func (s *SalesService) MoveDealStage(ctx context.Context, id string, stageID string) (*Deal, error) {
	var deal Deal
	body := map[string]string{"stage_id": stageID}
	err := s.client.request(ctx, "POST", "/sales/v1/deals/"+id+"/move", body, &deal)
	if err != nil {
		return nil, err
	}
	return &deal, nil
}

// ============================================================
// ACTIVITIES
// ============================================================

// Activity represents a sales activity such as a task, call, meeting, or email.
type Activity struct {
	ID              string     `json:"id"`
	Type            string     `json:"type"` // task, call, meeting, email
	Subject         string     `json:"subject"`
	Description     *string    `json:"description,omitempty"`
	DueDate         *string    `json:"due_date,omitempty"`
	DueTime         *string    `json:"due_time,omitempty"`
	DurationMinutes *int       `json:"duration_minutes,omitempty"`
	AssignedToID    *string    `json:"assigned_to_id,omitempty"`
	IsCompleted     bool       `json:"is_completed"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// CreateActivityParams are the parameters for creating an activity.
type CreateActivityParams struct {
	Type            string         `json:"type"` // task, call, meeting, email
	Subject         string         `json:"subject"`
	Description     *string        `json:"description,omitempty"`
	DueDate         *string        `json:"due_date,omitempty"`
	DueTime         *string        `json:"due_time,omitempty"`
	DurationMinutes *int           `json:"duration_minutes,omitempty"`
	AssignedToID    *string        `json:"assigned_to_id,omitempty"`
	Links           []ActivityLink `json:"links,omitempty"`
}

// CreateActivity creates a new activity.
func (s *SalesService) CreateActivity(ctx context.Context, params *CreateActivityParams) (*Activity, error) {
	var activity Activity
	err := s.client.request(ctx, "POST", "/sales/v1/activities", params, &activity)
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

// ListActivitiesParams are the parameters for listing activities.
type ListActivitiesParams struct {
	Type        string `json:"type,omitempty"`
	IsCompleted *bool  `json:"is_completed,omitempty"`
}

// ActivityList is a list of activities.
type ActivityList struct {
	Activities []Activity `json:"activities"`
	Total      int        `json:"total"`
}

// ListActivities lists activities with optional filters.
func (s *SalesService) ListActivities(ctx context.Context, params *ListActivitiesParams) (*ActivityList, error) {
	path := "/sales/v1/activities"
	if params != nil {
		query := ""
		if params.Type != "" {
			query += "type=" + params.Type
		}
		if params.IsCompleted != nil {
			if query != "" {
				query += "&"
			}
			query += fmt.Sprintf("is_completed=%t", *params.IsCompleted)
		}
		if query != "" {
			path += "?" + query
		}
	}
	var list ActivityList
	err := s.client.request(ctx, "GET", path, nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// GetActivity retrieves an activity by ID.
func (s *SalesService) GetActivity(ctx context.Context, id string) (*Activity, error) {
	var activity Activity
	err := s.client.request(ctx, "GET", "/sales/v1/activities/"+id, nil, &activity)
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

// UpdateActivityParams are the parameters for updating an activity.
type UpdateActivityParams struct {
	Type            *string `json:"type,omitempty"`
	Subject         *string `json:"subject,omitempty"`
	Description     *string `json:"description,omitempty"`
	DueDate         *string `json:"due_date,omitempty"`
	DueTime         *string `json:"due_time,omitempty"`
	DurationMinutes *int    `json:"duration_minutes,omitempty"`
}

// UpdateActivity updates an activity.
func (s *SalesService) UpdateActivity(ctx context.Context, id string, params *UpdateActivityParams) (*Activity, error) {
	var activity Activity
	err := s.client.request(ctx, "PATCH", "/sales/v1/activities/"+id, params, &activity)
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

// DeleteActivity deletes an activity.
func (s *SalesService) DeleteActivity(ctx context.Context, id string) error {
	return s.client.request(ctx, "DELETE", "/sales/v1/activities/"+id, nil, nil)
}

// CompleteActivity marks an activity as completed or uncompleted.
func (s *SalesService) CompleteActivity(ctx context.Context, id string, completed bool) (*Activity, error) {
	var activity Activity
	body := map[string]bool{"completed": completed}
	err := s.client.request(ctx, "POST", "/sales/v1/activities/"+id+"/complete", body, &activity)
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

// ============================================================
// NOTES
// ============================================================

// SalesNote represents a note attached to a sales entity.
type SalesNote struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SalesNoteLink represents a link between a note and an entity.
type SalesNoteLink struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	IsPrimary  bool   `json:"is_primary,omitempty"`
}

// CreateSalesNoteParams are the parameters for creating a sales note.
type CreateSalesNoteParams struct {
	Content  string          `json:"content"`
	AuthorID *string         `json:"author_id,omitempty"`
	Links    []SalesNoteLink `json:"links,omitempty"`
}

// CreateSalesNote creates a new sales note.
func (s *SalesService) CreateSalesNote(ctx context.Context, params *CreateSalesNoteParams) (*SalesNote, error) {
	var note SalesNote
	err := s.client.request(ctx, "POST", "/sales/v1/notes", params, &note)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// SalesNoteList is a list of sales notes.
type SalesNoteList struct {
	Notes []SalesNote `json:"notes"`
	Total int         `json:"total"`
}

// ListSalesNotes lists all sales notes in the workspace.
func (s *SalesService) ListSalesNotes(ctx context.Context) (*SalesNoteList, error) {
	var list SalesNoteList
	err := s.client.request(ctx, "GET", "/sales/v1/notes", nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// GetSalesNote retrieves a sales note by ID.
func (s *SalesService) GetSalesNote(ctx context.Context, id string) (*SalesNote, error) {
	var note SalesNote
	err := s.client.request(ctx, "GET", "/sales/v1/notes/"+id, nil, &note)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// UpdateSalesNoteParams are the parameters for updating a sales note.
type UpdateSalesNoteParams struct {
	Content *string `json:"content,omitempty"`
}

// UpdateSalesNote updates a sales note.
func (s *SalesService) UpdateSalesNote(ctx context.Context, id string, params *UpdateSalesNoteParams) (*SalesNote, error) {
	var note SalesNote
	err := s.client.request(ctx, "PATCH", "/sales/v1/notes/"+id, params, &note)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// DeleteSalesNote deletes a sales note.
func (s *SalesService) DeleteSalesNote(ctx context.Context, id string) error {
	return s.client.request(ctx, "DELETE", "/sales/v1/notes/"+id, nil, nil)
}

// ============================================================
// CONTACT BASES
// ============================================================

// ContactBase represents a contact base for organizing persons and organizations.
type ContactBase struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ContactBaseList is a list of contact bases.
type ContactBaseList struct {
	ContactBases []ContactBase `json:"contact_bases"`
	Total        int           `json:"total"`
}

// CreateContactBaseParams contains parameters for creating a contact base.
type CreateContactBaseParams struct {
	Name string `json:"name"`
}

// CreateContactBase creates a new contact base.
func (s *SalesService) CreateContactBase(ctx context.Context, params *CreateContactBaseParams) (*ContactBase, error) {
	var cb ContactBase
	err := s.client.request(ctx, "POST", "/sales/v1/contact-bases", params, &cb)
	if err != nil {
		return nil, err
	}
	return &cb, nil
}

// ListContactBases lists all contact bases.
func (s *SalesService) ListContactBases(ctx context.Context) (*ContactBaseList, error) {
	var list ContactBaseList
	err := s.client.request(ctx, "GET", "/sales/v1/contact-bases", nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// GetContactBase retrieves a contact base by ID.
func (s *SalesService) GetContactBase(ctx context.Context, id string) (*ContactBase, error) {
	var cb ContactBase
	err := s.client.request(ctx, "GET", "/sales/v1/contact-bases/"+id, nil, &cb)
	if err != nil {
		return nil, err
	}
	return &cb, nil
}

// ============================================================
// PERSONS
// ============================================================

// Person represents a contact person in the sales system.
type Person struct {
	ID          string    `json:"id"`
	FullName    string    `json:"full_name"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	LinkedInURL *string   `json:"linkedin_url,omitempty"`
	OwnerID     *string   `json:"owner_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreatePersonParams are the parameters for creating a person.
type CreatePersonParams struct {
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       string  `json:"email,omitempty"`
	Phone       string  `json:"phone,omitempty"`
	LinkedInURL *string `json:"linkedin_url,omitempty"`
	OwnerID     *string `json:"owner_id,omitempty"`
}

// CreatePerson creates a new person in a contact base.
func (s *SalesService) CreatePerson(ctx context.Context, contactBaseID string, params *CreatePersonParams) (*Person, error) {
	var person Person
	err := s.client.request(ctx, "POST", "/sales/v1/contact-bases/"+contactBaseID+"/persons", params, &person)
	if err != nil {
		return nil, err
	}
	return &person, nil
}

// PersonList is a list of persons.
type PersonList struct {
	Persons []Person `json:"persons"`
	Total   int      `json:"total"`
}

// ListPersons lists all persons in a contact base.
func (s *SalesService) ListPersons(ctx context.Context, contactBaseID string) (*PersonList, error) {
	var list PersonList
	err := s.client.request(ctx, "GET", "/sales/v1/contact-bases/"+contactBaseID+"/persons", nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// GetPerson retrieves a person by ID.
func (s *SalesService) GetPerson(ctx context.Context, id string) (*Person, error) {
	var person Person
	err := s.client.request(ctx, "GET", "/sales/v1/persons/"+id, nil, &person)
	if err != nil {
		return nil, err
	}
	return &person, nil
}

// UpdatePersonParams are the parameters for updating a person.
type UpdatePersonParams struct {
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	LinkedInURL *string `json:"linkedin_url,omitempty"`
	OwnerID     *string `json:"owner_id,omitempty"`
}

// UpdatePerson updates a person.
func (s *SalesService) UpdatePerson(ctx context.Context, id string, params *UpdatePersonParams) (*Person, error) {
	var person Person
	err := s.client.request(ctx, "PATCH", "/sales/v1/persons/"+id, params, &person)
	if err != nil {
		return nil, err
	}
	return &person, nil
}

// DeletePerson deletes a person.
func (s *SalesService) DeletePerson(ctx context.Context, id string) error {
	return s.client.request(ctx, "DELETE", "/sales/v1/persons/"+id, nil, nil)
}

// ============================================================
// ORGANIZATIONS
// ============================================================

// Organization represents a company or organization in the sales system.
type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Website     *string   `json:"website,omitempty"`
	LinkedInURL *string   `json:"linkedin_url,omitempty"`
	OwnerID     *string   `json:"owner_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateOrganizationParams are the parameters for creating an organization.
type CreateOrganizationParams struct {
	Name        string  `json:"name"`
	Website     *string `json:"website,omitempty"`
	LinkedInURL *string `json:"linkedin_url,omitempty"`
	OwnerID     *string `json:"owner_id,omitempty"`
}

// CreateOrganization creates a new organization in a contact base.
func (s *SalesService) CreateOrganization(ctx context.Context, contactBaseID string, params *CreateOrganizationParams) (*Organization, error) {
	var org Organization
	err := s.client.request(ctx, "POST", "/sales/v1/contact-bases/"+contactBaseID+"/organizations", params, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// OrganizationList is a list of organizations.
type OrganizationList struct {
	Organizations []Organization `json:"organizations"`
	Total         int            `json:"total"`
}

// ListOrganizations lists all organizations in a contact base.
func (s *SalesService) ListOrganizations(ctx context.Context, contactBaseID string) (*OrganizationList, error) {
	var list OrganizationList
	err := s.client.request(ctx, "GET", "/sales/v1/contact-bases/"+contactBaseID+"/organizations", nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// GetOrganization retrieves an organization by ID.
func (s *SalesService) GetOrganization(ctx context.Context, id string) (*Organization, error) {
	var org Organization
	err := s.client.request(ctx, "GET", "/sales/v1/organizations/"+id, nil, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// UpdateOrganizationParams are the parameters for updating an organization.
type UpdateOrganizationParams struct {
	Name        *string `json:"name,omitempty"`
	Website     *string `json:"website,omitempty"`
	LinkedInURL *string `json:"linkedin_url,omitempty"`
	OwnerID     *string `json:"owner_id,omitempty"`
}

// UpdateOrganization updates an organization.
func (s *SalesService) UpdateOrganization(ctx context.Context, id string, params *UpdateOrganizationParams) (*Organization, error) {
	var org Organization
	err := s.client.request(ctx, "PATCH", "/sales/v1/organizations/"+id, params, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// DeleteOrganization deletes an organization.
func (s *SalesService) DeleteOrganization(ctx context.Context, id string) error {
	return s.client.request(ctx, "DELETE", "/sales/v1/organizations/"+id, nil, nil)
}
