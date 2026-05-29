package tedo

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

const (
	ColumnTypeText           = "text"
	ColumnTypeNumber         = "number"
	ColumnTypeNumberWithUnit = "number_with_unit"
	ColumnTypeBoolean        = "boolean"
	ColumnTypeDate           = "date"
	ColumnTypeTimestamp      = "timestamp"
	ColumnTypeURL            = "url"
	ColumnTypeSelect         = "select"
	ColumnTypeMultiSelect    = "multi_select"
	ColumnTypeReference      = "reference"
	ColumnTypeMultiReference = "multi_reference"
	ColumnTypeRollup         = "rollup"
)

// TablesService handles Tables API calls.
type TablesService struct {
	client *Client
}

type JSONMap map[string]any

type TablesBase struct {
	ID                 string     `json:"id"`
	WorkspaceID        string     `json:"workspace_id"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	CreatedByActorType string     `json:"created_by_actor_type"`
	CreatedByActorID   string     `json:"created_by_actor_id"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	ArchivedAt         *time.Time `json:"archived_at"`
}

type Table struct {
	ID              string     `json:"id"`
	WorkspaceID     string     `json:"workspace_id"`
	BaseID          string     `json:"base_id"`
	Name            string     `json:"name"`
	Slug            string     `json:"slug"`
	Description     string     `json:"description"`
	PrimaryColumnID string     `json:"primary_column_id"`
	Position        int        `json:"position"`
	SchemaVersion   int        `json:"schema_version"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ArchivedAt      *time.Time `json:"archived_at"`
}

type Column struct {
	ID           string     `json:"id"`
	TableID      string     `json:"table_id"`
	Key          string     `json:"key"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	TypeConfig   JSONMap    `json:"type_config"`
	Required     bool       `json:"required"`
	DefaultValue any        `json:"default_value,omitempty"`
	Position     int        `json:"position"`
	Indexed      bool       `json:"indexed"`
	Unique       bool       `json:"unique"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ArchivedAt   *time.Time `json:"archived_at"`
}

type Row struct {
	ID                 string     `json:"id"`
	WorkspaceID        string     `json:"workspace_id"`
	TableID            string     `json:"table_id"`
	ExternalKey        string     `json:"external_key"`
	Data               JSONMap    `json:"data"`
	DataVersion        int        `json:"data_version"`
	CreatedByActorType string     `json:"created_by_actor_type"`
	CreatedByActorID   string     `json:"created_by_actor_id"`
	CreatedByActorName string     `json:"created_by_actor_name"`
	UpdatedByActorType string     `json:"updated_by_actor_type"`
	UpdatedByActorID   string     `json:"updated_by_actor_id"`
	UpdatedByActorName string     `json:"updated_by_actor_name"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	ArchivedAt         *time.Time `json:"archived_at"`
}

type HistoryEvent struct {
	ID              string    `json:"id"`
	TableID         string    `json:"table_id"`
	RowID           string    `json:"row_id"`
	OperationKey    string    `json:"operation_key"`
	Stage           string    `json:"stage"`
	Code            string    `json:"code"`
	ActorType       string    `json:"actor_type"`
	ActorID         string    `json:"actor_id"`
	ActorName       string    `json:"actor_name"`
	OriginActorType string    `json:"origin_actor_type"`
	OriginActorID   string    `json:"origin_actor_id"`
	SourceURL       string    `json:"source_url"`
	Patch           JSONMap   `json:"patch"`
	At              time.Time `json:"at"`
}

type Snapshot struct {
	ID                 string    `json:"id"`
	TableID            string    `json:"table_id"`
	Label              string    `json:"label"`
	SchemaFrozen       JSONMap   `json:"schema_frozen"`
	DataFrozenRef      string    `json:"data_frozen_ref"`
	CreatedByActorType string    `json:"created_by_actor_type"`
	CreatedByActorID   string    `json:"created_by_actor_id"`
	CreatedAt          time.Time `json:"created_at"`
}

type ShareToken struct {
	ID          string     `json:"id"`
	ScopeKind   string     `json:"scope_kind"`
	ScopeID     string     `json:"scope_id"`
	Permissions JSONMap    `json:"permissions"`
	Kind        string     `json:"kind"`
	Token       string     `json:"token"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
	RevokedAt   *time.Time `json:"revoked_at"`
}

type StatusResult struct {
	Status string `json:"status"`
}

type CreateBaseParams struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type UpdateBaseParams struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ListBasesParams struct {
	Limit int
}

type ColumnSeed struct {
	Name       string  `json:"name"`
	Key        string  `json:"key,omitempty"`
	Type       string  `json:"type"`
	TypeConfig JSONMap `json:"type_config,omitempty"`
	Required   bool    `json:"required,omitempty"`
	Position   int     `json:"position,omitempty"`
}

type CreateTableParams struct {
	BaseID      string       `json:"base_id,omitempty"`
	Name        string       `json:"name"`
	Slug        string       `json:"slug,omitempty"`
	Description string       `json:"description,omitempty"`
	Columns     []ColumnSeed `json:"columns,omitempty"`
}

type UpdateTableParams struct {
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

type ListTablesParams struct {
	BaseID string
	Limit  int
}

type TableWithSchema struct {
	Table  Table    `json:"table"`
	Schema []Column `json:"schema"`
}

type UpsertColumnParams struct {
	Name         string  `json:"name,omitempty"`
	Key          string  `json:"key,omitempty"`
	ColumnID     string  `json:"column_id,omitempty"`
	Type         string  `json:"type"`
	TypeConfig   JSONMap `json:"type_config,omitempty"`
	Required     *bool   `json:"required,omitempty"`
	DefaultValue any     `json:"default_value,omitempty"`
	Position     int     `json:"position,omitempty"`
	Indexed      *bool   `json:"indexed,omitempty"`
	Unique       *bool   `json:"unique,omitempty"`
}

type UpdateColumnParams struct {
	Name         string  `json:"name,omitempty"`
	Key          string  `json:"key,omitempty"`
	Type         string  `json:"type,omitempty"`
	TypeConfig   JSONMap `json:"type_config,omitempty"`
	Required     *bool   `json:"required,omitempty"`
	DefaultValue any     `json:"default_value,omitempty"`
	Position     int     `json:"position,omitempty"`
	Indexed      *bool   `json:"indexed,omitempty"`
	Unique       *bool   `json:"unique,omitempty"`
}

type UpsertColumnResult struct {
	Column   Column   `json:"column"`
	Warnings []string `json:"warnings,omitempty"`
}

type RetypeColumnParams struct {
	Type   string  `json:"type"`
	Config JSONMap `json:"config,omitempty"`
	Mode   string  `json:"mode,omitempty"`
}

type RetypeColumnResult struct {
	Column   Column   `json:"column"`
	Lossless bool     `json:"lossless"`
	Warnings []string `json:"warnings"`
}

type UpsertRowParams struct {
	Values    JSONMap `json:"values"`
	RowID     string  `json:"row_id,omitempty"`
	KeyColumn string  `json:"key_column,omitempty"`
	KeyValue  string  `json:"key_value,omitempty"`
}

type QueryRowsParams struct {
	Filter    string
	Fields    string
	Sort      string
	Expand    string
	V         string `json:"v,omitempty"`
	Cursor    string
	Limit     int
	GroupBy   string
	Aggregate string
	Page      int
}

type QueryRowsResult struct {
	Rows       []Row     `json:"rows,omitempty"`
	Schema     []Column  `json:"schema,omitempty"`
	NextCursor string    `json:"next_cursor,omitempty"`
	Grouped    bool      `json:"grouped,omitempty"`
	Groups     []JSONMap `json:"groups,omitempty"`
	HasMore    bool      `json:"has_more,omitempty"`
	NextPage   int       `json:"next_page,omitempty"`
	Meta       JSONMap   `json:"_meta,omitempty"`
}

type UpsertRowResult struct {
	Row     Row  `json:"row"`
	Created bool `json:"created"`
	Updated bool `json:"updated"`
}

type UpdateRowParams struct {
	Data JSONMap `json:"data"`
}

type BulkUpsertRowsParams struct {
	Rows          []JSONMap `json:"rows"`
	KeyColumn     string    `json:"key_column,omitempty"`
	DeleteMissing *bool     `json:"delete_missing,omitempty"`
	BulkReplace   *bool     `json:"-"`
}

type BulkUpsertRowsResult struct {
	Rows      []Row    `json:"rows"`
	Inserted  int      `json:"inserted"`
	Updated   int      `json:"updated"`
	Deleted   int      `json:"deleted"`
	Unchanged int      `json:"unchanged"`
	Failed    int      `json:"failed"`
	Warnings  []string `json:"warnings"`
	Created   int      `json:"created"`
}

type CreateSnapshotParams struct {
	Label string `json:"label,omitempty"`
}

type CreateShareTokenParams struct {
	Kind        string  `json:"kind,omitempty"`
	Permissions JSONMap `json:"permissions,omitempty"`
	ExpiresAt   string  `json:"expires_at,omitempty"`
}

type CreateShareTokenResult struct {
	ShareToken ShareToken `json:"share_token"`
	PublicURL  string     `json:"public_url"`
}

type ImportCSVParams struct {
	BaseID       string            `json:"base_id,omitempty"`
	TableID      string            `json:"table_id,omitempty"`
	Name         string            `json:"name,omitempty"`
	CSVBody      string            `json:"csv_body"`
	KeyColumn    string            `json:"key_column"`
	Mode         string            `json:"mode,omitempty"`
	ColumnTypes  map[string]string `json:"column_types,omitempty"`
	UnitDefaults map[string]string `json:"unit_defaults,omitempty"`
}

type ImportCSVResult struct {
	Import   JSONMap  `json:"import,omitempty"`
	Table    Table    `json:"table,omitempty"`
	Counts   JSONMap  `json:"counts,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

func (s *TablesService) CreateBase(ctx context.Context, params *CreateBaseParams) (*TablesBase, error) {
	var out struct {
		Base TablesBase `json:"base"`
	}
	if err := s.client.request(ctx, "POST", "/tables/v1/bases", params, &out); err != nil {
		return nil, err
	}
	return &out.Base, nil
}

func (s *TablesService) ListBases(ctx context.Context, params *ListBasesParams) ([]TablesBase, error) {
	var out struct {
		Bases []TablesBase `json:"bases"`
	}
	if err := s.client.request(ctx, "GET", "/tables/v1/bases"+queryString(params), nil, &out); err != nil {
		return nil, err
	}
	return out.Bases, nil
}

func (s *TablesService) GetBase(ctx context.Context, baseID string) (*TablesBase, error) {
	var out struct {
		Base TablesBase `json:"base"`
	}
	if err := s.client.request(ctx, "GET", "/tables/v1/bases/"+pathEscape(baseID), nil, &out); err != nil {
		return nil, err
	}
	return &out.Base, nil
}

func (s *TablesService) UpdateBase(ctx context.Context, baseID string, params *UpdateBaseParams) (*TablesBase, error) {
	var out struct {
		Base TablesBase `json:"base"`
	}
	if err := s.client.request(ctx, "PATCH", "/tables/v1/bases/"+pathEscape(baseID), params, &out); err != nil {
		return nil, err
	}
	return &out.Base, nil
}

func (s *TablesService) ArchiveBase(ctx context.Context, baseID string) (*StatusResult, error) {
	var out StatusResult
	if err := s.client.request(ctx, "DELETE", "/tables/v1/bases/"+pathEscape(baseID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TablesService) CreateTable(ctx context.Context, params *CreateTableParams) (*TableWithSchema, error) {
	var out TableWithSchema
	if err := s.client.request(ctx, "POST", "/tables/v1/tables", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TablesService) CreateTableInBase(ctx context.Context, baseID string, params *CreateTableParams) (*TableWithSchema, error) {
	var out TableWithSchema
	if err := s.client.request(ctx, "POST", "/tables/v1/bases/"+pathEscape(baseID)+"/tables", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TablesService) ListTables(ctx context.Context, params *ListTablesParams) ([]Table, error) {
	var out struct {
		Tables []Table `json:"tables"`
	}
	if err := s.client.request(ctx, "GET", "/tables/v1/tables"+queryString(params), nil, &out); err != nil {
		return nil, err
	}
	return out.Tables, nil
}

func (s *TablesService) GetTable(ctx context.Context, tableID string) (*TableWithSchema, error) {
	var out TableWithSchema
	if err := s.client.request(ctx, "GET", "/tables/v1/tables/"+pathEscape(tableID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TablesService) UpdateTable(ctx context.Context, tableID string, params *UpdateTableParams) (*Table, error) {
	var out struct {
		Table Table `json:"table"`
	}
	if err := s.client.request(ctx, "PATCH", "/tables/v1/tables/"+pathEscape(tableID), params, &out); err != nil {
		return nil, err
	}
	return &out.Table, nil
}

func (s *TablesService) ArchiveTable(ctx context.Context, tableID string) (*StatusResult, error) {
	var out StatusResult
	if err := s.client.request(ctx, "DELETE", "/tables/v1/tables/"+pathEscape(tableID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TablesService) UpsertColumn(ctx context.Context, tableID string, params *UpsertColumnParams) (*UpsertColumnResult, error) {
	var out UpsertColumnResult
	if err := s.client.request(ctx, "POST", "/tables/v1/tables/"+pathEscape(tableID)+"/columns", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TablesService) ListColumns(ctx context.Context, tableID string) ([]Column, error) {
	var out struct {
		Columns []Column `json:"columns"`
	}
	if err := s.client.request(ctx, "GET", "/tables/v1/tables/"+pathEscape(tableID)+"/columns", nil, &out); err != nil {
		return nil, err
	}
	return out.Columns, nil
}

func (s *TablesService) UpdateColumn(ctx context.Context, columnID string, params *UpdateColumnParams) (*Column, error) {
	var out struct {
		Column Column `json:"column"`
	}
	if err := s.client.request(ctx, "PATCH", "/tables/v1/columns/"+pathEscape(columnID), params, &out); err != nil {
		return nil, err
	}
	return &out.Column, nil
}

func (s *TablesService) ArchiveColumn(ctx context.Context, columnID string) (*StatusResult, error) {
	var out StatusResult
	if err := s.client.request(ctx, "DELETE", "/tables/v1/columns/"+pathEscape(columnID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TablesService) RetypeColumn(ctx context.Context, tableID, columnID string, params *RetypeColumnParams) (*RetypeColumnResult, error) {
	var out RetypeColumnResult
	if err := s.client.request(ctx, "POST", "/tables/v1/tables/"+pathEscape(tableID)+"/columns/"+pathEscape(columnID)+"/retype", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TablesService) UpsertRow(ctx context.Context, tableID string, params *UpsertRowParams) (*UpsertRowResult, error) {
	var out UpsertRowResult
	if err := s.client.request(ctx, "POST", "/tables/v1/tables/"+pathEscape(tableID)+"/rows", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TablesService) QueryRows(ctx context.Context, tableID string, params *QueryRowsParams) (*QueryRowsResult, error) {
	var out QueryRowsResult
	if err := s.client.request(ctx, "GET", "/tables/v1/tables/"+pathEscape(tableID)+"/rows"+queryString(params), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TablesService) GetRow(ctx context.Context, tableID, rowID string) (*Row, error) {
	var out struct {
		Row Row `json:"row"`
	}
	if err := s.client.request(ctx, "GET", "/tables/v1/tables/"+pathEscape(tableID)+"/rows/"+pathEscape(rowID), nil, &out); err != nil {
		return nil, err
	}
	return &out.Row, nil
}

func (s *TablesService) RowHistory(ctx context.Context, tableID, rowID string) ([]HistoryEvent, error) {
	var out struct {
		Events []HistoryEvent `json:"events"`
	}
	if err := s.client.request(ctx, "GET", "/tables/v1/tables/"+pathEscape(tableID)+"/rows/"+pathEscape(rowID)+"/history", nil, &out); err != nil {
		return nil, err
	}
	return out.Events, nil
}

func (s *TablesService) UpdateRow(ctx context.Context, tableID, rowID string, params *UpdateRowParams) (*Row, error) {
	var out struct {
		Row Row `json:"row"`
	}
	if err := s.client.request(ctx, "PATCH", "/tables/v1/tables/"+pathEscape(tableID)+"/rows/"+pathEscape(rowID), params, &out); err != nil {
		return nil, err
	}
	return &out.Row, nil
}

func (s *TablesService) DeleteRow(ctx context.Context, tableID, rowID string) (*StatusResult, error) {
	var out StatusResult
	if err := s.client.request(ctx, "DELETE", "/tables/v1/tables/"+pathEscape(tableID)+"/rows/"+pathEscape(rowID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TablesService) BulkUpsertRows(ctx context.Context, tableID string, params *BulkUpsertRowsParams) (*BulkUpsertRowsResult, error) {
	if params != nil && params.DeleteMissing == nil && params.BulkReplace != nil {
		params.DeleteMissing = params.BulkReplace
	}
	var out BulkUpsertRowsResult
	if err := s.client.request(ctx, "POST", "/tables/v1/tables/"+pathEscape(tableID)+"/rows/bulk_upsert", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TablesService) CreateSnapshot(ctx context.Context, tableID string, params *CreateSnapshotParams) (*Snapshot, error) {
	var out struct {
		Snapshot Snapshot `json:"snapshot"`
	}
	if err := s.client.request(ctx, "POST", "/tables/v1/tables/"+pathEscape(tableID)+"/snapshots", params, &out); err != nil {
		return nil, err
	}
	return &out.Snapshot, nil
}

func (s *TablesService) ListSnapshots(ctx context.Context, tableID string) ([]Snapshot, error) {
	var out struct {
		Snapshots []Snapshot `json:"snapshots"`
	}
	if err := s.client.request(ctx, "GET", "/tables/v1/tables/"+pathEscape(tableID)+"/snapshots", nil, &out); err != nil {
		return nil, err
	}
	return out.Snapshots, nil
}

func (s *TablesService) RestoreSnapshot(ctx context.Context, tableID, snapshotID string) (*Snapshot, error) {
	var out struct {
		Snapshot Snapshot `json:"snapshot"`
	}
	if err := s.client.request(ctx, "POST", "/tables/v1/tables/"+pathEscape(tableID)+"/snapshots/"+pathEscape(snapshotID)+"/restore", nil, &out); err != nil {
		return nil, err
	}
	return &out.Snapshot, nil
}

func (s *TablesService) CreateShareToken(ctx context.Context, tableID string, params *CreateShareTokenParams) (*CreateShareTokenResult, error) {
	var out CreateShareTokenResult
	if err := s.client.request(ctx, "POST", "/tables/v1/tables/"+pathEscape(tableID)+"/shares", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TablesService) RevokeShareToken(ctx context.Context, tokenID string) (*StatusResult, error) {
	var out StatusResult
	if err := s.client.request(ctx, "DELETE", "/tables/v1/share_tokens/"+pathEscape(tokenID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *TablesService) ImportCSV(ctx context.Context, params *ImportCSVParams) (*ImportCSVResult, error) {
	var out ImportCSVResult
	if err := s.client.request(ctx, "POST", "/tables/v1/import/csv", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func queryString(params any) string {
	values := url.Values{}
	switch p := params.(type) {
	case *ListBasesParams:
		if p != nil && p.Limit > 0 {
			values.Set("limit", strconv.Itoa(p.Limit))
		}
	case *ListTablesParams:
		if p != nil {
			if p.BaseID != "" {
				values.Set("base_id", p.BaseID)
			}
			if p.Limit > 0 {
				values.Set("limit", strconv.Itoa(p.Limit))
			}
		}
	case *QueryRowsParams:
		if p != nil {
			setString(values, "filter", p.Filter)
			setString(values, "fields", p.Fields)
			setString(values, "sort", p.Sort)
			setString(values, "expand", p.Expand)
			setString(values, "v", p.V)
			setString(values, "cursor", p.Cursor)
			setString(values, "group_by", p.GroupBy)
			setString(values, "aggregate", p.Aggregate)
			if p.Limit > 0 {
				values.Set("limit", strconv.Itoa(p.Limit))
			}
			if p.Page > 0 {
				values.Set("page", strconv.Itoa(p.Page))
			}
		}
	}
	if len(values) == 0 {
		return ""
	}
	return "?" + values.Encode()
}

func setString(values url.Values, key, value string) {
	if value != "" {
		values.Set(key, value)
	}
}

func pathEscape(value string) string {
	return url.PathEscape(value)
}
