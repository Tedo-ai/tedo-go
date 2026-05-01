package tedo

import "time"

// Invoice represents an issued invoice in the Invoicing app.
//
// Integer-cent fields (SubtotalCents, TaxCents, TotalCents) are canonical
// after Phase 1.1. The legacy float fields (SubtotalAmount, TaxAmount,
// TotalAmount) are retained for the migration window and will be removed
// once the follow-up drop migration ships.
type Invoice struct {
	ID            string         `json:"id"`
	WorkspaceID   string         `json:"workspace_id"`
	LegalEntityID string         `json:"legal_entity_id"`
	RelationID    *string        `json:"relation_id,omitempty"`
	SequenceID    *string        `json:"sequence_id,omitempty"`
	InvoiceNumber *string        `json:"invoice_number,omitempty"`
	Status        string         `json:"status"`
	IssueDate     *time.Time     `json:"issue_date,omitempty"`
	DueDate       *time.Time     `json:"due_date,omitempty"`
	Currency      string         `json:"currency"`
	SubtotalCents int64          `json:"subtotal_cents"`
	TaxCents      int64          `json:"tax_cents"`
	TotalCents    int64          `json:"total_cents"`
	SubtotalAmount *float64      `json:"subtotal_amount,omitempty"` // deprecated
	TaxAmount      *float64      `json:"tax_amount,omitempty"`      // deprecated
	TotalAmount    *float64      `json:"total_amount,omitempty"`    // deprecated
	Notes         *string        `json:"notes,omitempty"`
	TaxReason     *string        `json:"tax_reason,omitempty"`
	Lines         []InvoiceLine  `json:"lines,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// InvoiceLine represents a single line item on an Invoice.
//
// Integer-cent fields (UnitPriceCents, LineTotalCents, TaxTotalCents) are
// canonical after Phase 1.1. TaxPercentageBP stores the rate in basis points
// (21% → 2100). Legacy float fields are retained for the migration window.
type InvoiceLine struct {
	ID              string   `json:"id"`
	InvoiceID       string   `json:"invoice_id"`
	Description     string   `json:"description"`
	Quantity        float64  `json:"quantity"`
	UnitPriceCents  int64    `json:"unit_price_cents"`
	LineTotalCents  int64    `json:"line_total_cents"`
	TaxTotalCents   int64    `json:"tax_total_cents"`
	TaxPercentageBP int32    `json:"tax_percentage_bp"` // basis points: 21% -> 2100
	UnitPrice       *float64 `json:"unit_price,omitempty"`  // deprecated
	LineTotal       *float64 `json:"line_total,omitempty"`  // deprecated
	TaxTotal        *float64 `json:"tax_total,omitempty"`   // deprecated
	TaxPercentage   *float64 `json:"tax_percentage,omitempty"` // deprecated
	SortOrder       int      `json:"sort_order"`
}
