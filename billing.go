package tedo

import (
	"context"
	"fmt"
	"time"
)

// BillingService handles billing-related API calls.
type BillingService struct {
	client *Client
}

// ============================================================
// CUSTOMERS
// ============================================================

// Customer represents a billing customer.
type Customer struct {
	ID            string            `json:"id"`
	Email         string            `json:"email"`
	Name          string            `json:"name,omitempty"`
	ExternalID    string            `json:"external_id,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	Subscriptions []Subscription    `json:"subscriptions,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at,omitempty"`
}

// CreateCustomerParams are the parameters for creating a customer.
type CreateCustomerParams struct {
	Email      string            `json:"email"`
	Name       string            `json:"name,omitempty"`
	ExternalID string            `json:"external_id,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// CreateCustomer creates a new customer.
func (s *BillingService) CreateCustomer(ctx context.Context, params *CreateCustomerParams) (*Customer, error) {
	var customer Customer
	err := s.client.request(ctx, "POST", "/billing/customers", params, &customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// GetCustomer retrieves a customer by ID.
func (s *BillingService) GetCustomer(ctx context.Context, id string) (*Customer, error) {
	var customer Customer
	err := s.client.request(ctx, "GET", "/billing/customers/"+id, nil, &customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// ListCustomersParams are the parameters for listing customers.
type ListCustomersParams struct {
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// CustomerList is a paginated list of customers.
type CustomerList struct {
	Customers  []Customer `json:"customers"`
	Total      int        `json:"total"`
	NextCursor string     `json:"next_cursor,omitempty"`
}

// ListCustomers lists all customers.
func (s *BillingService) ListCustomers(ctx context.Context, params *ListCustomersParams) (*CustomerList, error) {
	path := "/billing/customers"
	if params != nil {
		query := ""
		if params.Limit > 0 {
			query += fmt.Sprintf("limit=%d", params.Limit)
		}
		if params.Cursor != "" {
			if query != "" {
				query += "&"
			}
			query += "cursor=" + params.Cursor
		}
		if query != "" {
			path += "?" + query
		}
	}

	var list CustomerList
	err := s.client.request(ctx, "GET", path, nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// UpdateCustomerParams are the parameters for updating a customer.
type UpdateCustomerParams struct {
	Email      *string            `json:"email,omitempty"`
	Name       *string            `json:"name,omitempty"`
	ExternalID *string            `json:"external_id,omitempty"`
	Metadata   map[string]string  `json:"metadata,omitempty"`
}

// UpdateCustomer updates a customer.
func (s *BillingService) UpdateCustomer(ctx context.Context, id string, params *UpdateCustomerParams) (*Customer, error) {
	var customer Customer
	err := s.client.request(ctx, "PATCH", "/billing/customers/"+id, params, &customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// DeleteCustomer deletes a customer.
func (s *BillingService) DeleteCustomer(ctx context.Context, id string) error {
	return s.client.request(ctx, "DELETE", "/billing/customers/"+id, nil, nil)
}

// ============================================================
// SUBSCRIPTIONS
// ============================================================

// Subscription represents a billing subscription.
type Subscription struct {
	ID         string     `json:"id"`
	CustomerID string     `json:"customer_id"`
	PriceID    string     `json:"price_id"`
	Status     string     `json:"status"` // active, canceled, past_due
	Quantity   int        `json:"quantity,omitempty"`
	StartedAt  time.Time  `json:"started_at"`
	CanceledAt *time.Time `json:"canceled_at,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// CreateSubscriptionParams are the parameters for creating a subscription.
type CreateSubscriptionParams struct {
	CustomerID string            `json:"customer_id"`
	PriceID    string            `json:"price_id"`
	Quantity   int               `json:"quantity,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// CreateSubscription creates a new subscription.
func (s *BillingService) CreateSubscription(ctx context.Context, params *CreateSubscriptionParams) (*Subscription, error) {
	var subscription Subscription
	err := s.client.request(ctx, "POST", "/billing/subscriptions", params, &subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// GetSubscription retrieves a subscription by ID.
func (s *BillingService) GetSubscription(ctx context.Context, id string) (*Subscription, error) {
	var subscription Subscription
	err := s.client.request(ctx, "GET", "/billing/subscriptions/"+id, nil, &subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// CancelSubscription cancels a subscription.
func (s *BillingService) CancelSubscription(ctx context.Context, id string) (*Subscription, error) {
	var subscription Subscription
	err := s.client.request(ctx, "DELETE", "/billing/subscriptions/"+id, nil, &subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// ============================================================
// ENTITLEMENTS
// ============================================================

// EntitlementCheck is the result of an entitlement check.
type EntitlementCheck struct {
	HasAccess bool   `json:"has_access"`
	Value     string `json:"value,omitempty"`
	PlanName  string `json:"plan_name,omitempty"`
}

// CheckEntitlementParams are the parameters for checking an entitlement.
type CheckEntitlementParams struct {
	CustomerID     string `json:"customer_id"`
	EntitlementKey string `json:"entitlement_key"`
}

// CheckEntitlement checks if a customer has access to a feature.
func (s *BillingService) CheckEntitlement(ctx context.Context, params *CheckEntitlementParams) (*EntitlementCheck, error) {
	var result EntitlementCheck
	err := s.client.request(ctx, "POST", "/billing/entitlements/check", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ============================================================
// USAGE
// ============================================================

// UsageRecord represents a recorded usage event.
type UsageRecord struct {
	ID             string    `json:"id"`
	SubscriptionID string    `json:"subscription_id"`
	Quantity       int       `json:"quantity"`
	Timestamp      time.Time `json:"timestamp"`
}

// RecordUsageParams are the parameters for recording usage.
type RecordUsageParams struct {
	SubscriptionID string     `json:"subscription_id"`
	Quantity       int        `json:"quantity"`
	Timestamp      *time.Time `json:"timestamp,omitempty"`
	IdempotencyKey string     `json:"idempotency_key,omitempty"`
}

// RecordUsage records usage for a metered subscription.
func (s *BillingService) RecordUsage(ctx context.Context, params *RecordUsageParams) (*UsageRecord, error) {
	var record UsageRecord
	err := s.client.request(ctx, "POST", "/billing/usage", params, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// UsageSummary is an aggregated usage summary.
type UsageSummary struct {
	SubscriptionID string    `json:"subscription_id"`
	PeriodStart    time.Time `json:"period_start"`
	PeriodEnd      time.Time `json:"period_end"`
	TotalUsage     int       `json:"total_usage"`
	Records        int       `json:"records"`
}

// GetUsageSummaryParams are the parameters for getting a usage summary.
type GetUsageSummaryParams struct {
	SubscriptionID string `json:"subscription_id"`
}

// GetUsageSummary gets aggregated usage for a subscription.
func (s *BillingService) GetUsageSummary(ctx context.Context, params *GetUsageSummaryParams) (*UsageSummary, error) {
	path := "/billing/usage?subscription_id=" + params.SubscriptionID

	var summary UsageSummary
	err := s.client.request(ctx, "GET", path, nil, &summary)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}
