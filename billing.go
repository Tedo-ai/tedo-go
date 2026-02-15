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

// Plan/price keys for Tedo's built-in billing plans.
const (
	GuestPlanKey  = "guest"
	GuestPriceKey = "guest_monthly_EUR"
	FreePlanKey   = "free"
	FreePriceKey  = "free_monthly_EUR"
	BasicPlanKey  = "basic"
	BasicPriceKey = "basic_monthly_EUR"
)

// ============================================================
// PLANS
// ============================================================

// Plan represents a subscription plan.
type Plan struct {
	ID           string        `json:"id"`
	Key          string        `json:"key"`
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	IsActive     bool          `json:"is_active"`
	Prices       []Price       `json:"prices,omitempty"`
	Entitlements []Entitlement `json:"entitlements,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at,omitempty"`
}

// CreatePlanParams are the parameters for creating a plan.
type CreatePlanParams struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// CreatePlan creates a new subscription plan.
func (s *BillingService) CreatePlan(ctx context.Context, params *CreatePlanParams) (*Plan, error) {
	var plan Plan
	err := s.client.request(ctx, "POST", "/billing/v1/plans", params, &plan)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// PlanList is a list of plans.
type PlanList struct {
	Plans []Plan `json:"plans"`
	Total int    `json:"total"`
}

// ListPlans lists all plans.
func (s *BillingService) ListPlans(ctx context.Context) (*PlanList, error) {
	var list PlanList
	err := s.client.request(ctx, "GET", "/billing/v1/plans", nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// GetPlan retrieves a plan by ID.
func (s *BillingService) GetPlan(ctx context.Context, id string) (*Plan, error) {
	var plan Plan
	err := s.client.request(ctx, "GET", "/billing/v1/plans/"+id, nil, &plan)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// UpdatePlanParams are the parameters for updating a plan.
type UpdatePlanParams struct {
	Key         *string `json:"key,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// UpdatePlan updates a plan.
func (s *BillingService) UpdatePlan(ctx context.Context, id string, params *UpdatePlanParams) (*Plan, error) {
	var plan Plan
	err := s.client.request(ctx, "PATCH", "/billing/v1/plans/"+id, params, &plan)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// DeletePlan deletes (deactivates) a plan.
func (s *BillingService) DeletePlan(ctx context.Context, id string) error {
	return s.client.request(ctx, "DELETE", "/billing/v1/plans/"+id, nil, nil)
}

// ============================================================
// PRICES
// ============================================================

// Price represents a price for a plan.
type Price struct {
	ID            string    `json:"id"`
	PlanID        string    `json:"plan_id"`
	Key           string    `json:"key"`
	Amount        int       `json:"amount"` // in cents
	Currency      string    `json:"currency"`
	Interval      string    `json:"interval"` // month, year
	IntervalCount int       `json:"interval_count"`
	TrialDays     int       `json:"trial_days,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreatePriceParams are the parameters for creating a price.
type CreatePriceParams struct {
	Key           string `json:"key"`
	Amount        int    `json:"amount"`
	Currency      string `json:"currency,omitempty"`
	Interval      string `json:"interval,omitempty"`
	IntervalCount int    `json:"interval_count,omitempty"`
	TrialDays     int    `json:"trial_days,omitempty"`
}

// CreatePrice creates a new price for a plan.
func (s *BillingService) CreatePrice(ctx context.Context, planID string, params *CreatePriceParams) (*Price, error) {
	var price Price
	err := s.client.request(ctx, "POST", "/billing/v1/plans/"+planID+"/prices", params, &price)
	if err != nil {
		return nil, err
	}
	return &price, nil
}

// PriceList is a list of prices.
type PriceList struct {
	Prices []Price `json:"prices"`
	Total  int     `json:"total"`
}

// ListPrices lists all prices for a plan.
func (s *BillingService) ListPrices(ctx context.Context, planID string) (*PriceList, error) {
	var list PriceList
	err := s.client.request(ctx, "GET", "/billing/v1/plans/"+planID+"/prices", nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// ArchivePrice archives a price.
func (s *BillingService) ArchivePrice(ctx context.Context, planID, priceID string) error {
	return s.client.request(ctx, "DELETE", "/billing/v1/plans/"+planID+"/prices/"+priceID, nil, nil)
}

// ============================================================
// ENTITLEMENTS (Plan Features)
// ============================================================

// Entitlement represents a feature/limit on a plan.
type Entitlement struct {
	ID           string    `json:"id"`
	PlanID       string    `json:"plan_id"`
	Key          string    `json:"key"`
	ValueBool    *bool     `json:"value_bool,omitempty"`
	ValueInt     *int      `json:"value_int,omitempty"`
	OveragePrice int       `json:"overage_price,omitempty"`
	OverageUnit  int       `json:"overage_unit,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// CreateEntitlementParams are the parameters for creating an entitlement.
type CreateEntitlementParams struct {
	Key          string `json:"key"`
	ValueBool    *bool  `json:"value_bool,omitempty"`
	ValueInt     *int   `json:"value_int,omitempty"`
	OveragePrice int    `json:"overage_price,omitempty"`
	OverageUnit  int    `json:"overage_unit,omitempty"`
}

// CreateEntitlement creates an entitlement for a plan.
func (s *BillingService) CreateEntitlement(ctx context.Context, planID string, params *CreateEntitlementParams) (*Entitlement, error) {
	var entitlement Entitlement
	err := s.client.request(ctx, "POST", "/billing/v1/plans/"+planID+"/entitlements", params, &entitlement)
	if err != nil {
		return nil, err
	}
	return &entitlement, nil
}

// EntitlementList is a list of entitlements.
type EntitlementList struct {
	Entitlements []Entitlement `json:"entitlements"`
	Total        int           `json:"total"`
}

// ListEntitlements lists all entitlements for a plan.
func (s *BillingService) ListEntitlements(ctx context.Context, planID string) (*EntitlementList, error) {
	var list EntitlementList
	err := s.client.request(ctx, "GET", "/billing/v1/plans/"+planID+"/entitlements", nil, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// ArchiveEntitlement archives an entitlement.
func (s *BillingService) ArchiveEntitlement(ctx context.Context, planID, entitlementID string) error {
	return s.client.request(ctx, "DELETE", "/billing/v1/plans/"+planID+"/entitlements/"+entitlementID, nil, nil)
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
	err := s.client.request(ctx, "POST", "/billing/v1/customers", params, &customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// CreateCustomerForUser creates a billing customer for a user.
// The customer's ExternalID is set to "user:{userID}" for cross-referencing.
// Returns the customer ID.
func (s *BillingService) CreateCustomerForUser(ctx context.Context, userID int, email, name string) (string, error) {
	customer, err := s.CreateCustomer(ctx, &CreateCustomerParams{
		Email:      email,
		Name:       name,
		ExternalID: fmt.Sprintf("user:%d", userID),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create billing customer for user: %w", err)
	}
	return customer.ID, nil
}

// GetCustomer retrieves a customer by ID.
func (s *BillingService) GetCustomer(ctx context.Context, id string) (*Customer, error) {
	var customer Customer
	err := s.client.request(ctx, "GET", "/billing/v1/customers/"+id, nil, &customer)
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
	path := "/billing/v1/customers"
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
	Email      *string           `json:"email,omitempty"`
	Name       *string           `json:"name,omitempty"`
	ExternalID *string           `json:"external_id,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// UpdateCustomer updates a customer.
func (s *BillingService) UpdateCustomer(ctx context.Context, id string, params *UpdateCustomerParams) (*Customer, error) {
	var customer Customer
	err := s.client.request(ctx, "PATCH", "/billing/v1/customers/"+id, params, &customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// DeleteCustomer deletes a customer.
func (s *BillingService) DeleteCustomer(ctx context.Context, id string) error {
	return s.client.request(ctx, "DELETE", "/billing/v1/customers/"+id, nil, nil)
}

// ============================================================
// SUBSCRIPTIONS
// ============================================================

// Subscription represents a billing subscription.
type Subscription struct {
	ID         string            `json:"id"`
	CustomerID string            `json:"customer_id"`
	PriceID    string            `json:"price_id"`
	Status     string            `json:"status"` // active, canceled, past_due
	Quantity   int               `json:"quantity,omitempty"`
	StartedAt  time.Time         `json:"started_at"`
	CanceledAt *time.Time        `json:"canceled_at,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
}

// CreateSubscriptionParams are the parameters for creating a subscription.
type CreateSubscriptionParams struct {
	CustomerID string            `json:"customer_id"`
	PriceID    string            `json:"price_id,omitempty"`
	PlanKey    string            `json:"plan_key,omitempty"`
	PriceKey   string            `json:"price_key,omitempty"`
	Quantity   int               `json:"quantity,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// CreateSubscription creates a new subscription.
func (s *BillingService) CreateSubscription(ctx context.Context, params *CreateSubscriptionParams) (*Subscription, error) {
	var subscription Subscription
	err := s.client.request(ctx, "POST", "/billing/v1/subscriptions", params, &subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// CreateSubscriptionForWorkspace creates a free-tier subscription for a workspace.
// Returns the subscription ID.
func (s *BillingService) CreateSubscriptionForWorkspace(ctx context.Context, customerID, workspaceID string) (string, error) {
	return s.createSubscriptionWithPlan(ctx, customerID, FreePlanKey, FreePriceKey)
}

// CreateSubscriptionForGuestWorkspace creates a guest-tier subscription (lower limits).
// Returns the subscription ID.
func (s *BillingService) CreateSubscriptionForGuestWorkspace(ctx context.Context, customerID, workspaceID string) (string, error) {
	return s.createSubscriptionWithPlan(ctx, customerID, GuestPlanKey, GuestPriceKey)
}

// CreateSubscriptionForBasicPlan creates a basic paid subscription.
// Returns the subscription ID.
func (s *BillingService) CreateSubscriptionForBasicPlan(ctx context.Context, customerID string) (string, error) {
	return s.createSubscriptionWithPlan(ctx, customerID, BasicPlanKey, BasicPriceKey)
}

func (s *BillingService) createSubscriptionWithPlan(ctx context.Context, customerID, planKey, priceKey string) (string, error) {
	subscription, err := s.CreateSubscription(ctx, &CreateSubscriptionParams{
		CustomerID: customerID,
		PlanKey:    planKey,
		PriceKey:   priceKey,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create subscription: %w", err)
	}
	return subscription.ID, nil
}

// GetSubscription retrieves a subscription by ID.
func (s *BillingService) GetSubscription(ctx context.Context, id string) (*Subscription, error) {
	var subscription Subscription
	err := s.client.request(ctx, "GET", "/billing/v1/subscriptions/"+id, nil, &subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// CancelSubscription cancels a subscription.
func (s *BillingService) CancelSubscription(ctx context.Context, id string) (*Subscription, error) {
	var subscription Subscription
	err := s.client.request(ctx, "DELETE", "/billing/v1/subscriptions/"+id, nil, &subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// ============================================================
// CHECKOUT
// ============================================================

// CheckoutLink represents a billing checkout link.
type CheckoutLink struct {
	CheckoutURL string    `json:"checkout_url"`
	Token       string    `json:"token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// CreateCheckoutLinkParams are the parameters for creating a checkout link.
type CreateCheckoutLinkParams struct {
	ExpiresInHours int `json:"expires_in_hours,omitempty"`
}

// CreateCheckoutLink generates a checkout link for a subscription.
func (s *BillingService) CreateCheckoutLink(ctx context.Context, subscriptionID string, params *CreateCheckoutLinkParams) (*CheckoutLink, error) {
	var link CheckoutLink
	err := s.client.request(ctx, "POST", "/billing/v1/subscriptions/"+subscriptionID+"/checkout-link", params, &link)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// ============================================================
// ENTITLEMENT CHECK
// ============================================================

// EntitlementCheck is the result of an entitlement check.
type EntitlementCheck struct {
	HasAccess bool   `json:"has_access"`
	Value     any    `json:"value,omitempty"`
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
	err := s.client.request(ctx, "POST", "/billing/v1/entitlements/check", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CheckEntitlementByKey is a convenience method that checks an entitlement by customer ID and key.
func (s *BillingService) CheckEntitlementByKey(ctx context.Context, customerID, entitlementKey string) (*EntitlementCheck, error) {
	return s.CheckEntitlement(ctx, &CheckEntitlementParams{
		CustomerID:     customerID,
		EntitlementKey: entitlementKey,
	})
}

// ============================================================
// USAGE
// ============================================================

// UsageRecord represents a recorded usage event.
type UsageRecord struct {
	ID             string    `json:"id"`
	CustomerID     string    `json:"customer_id,omitempty"`
	SubscriptionID string    `json:"subscription_id,omitempty"`
	ProductKey     string    `json:"product_key"`
	Quantity       int       `json:"quantity"`
	Timestamp      time.Time `json:"timestamp"`
	IdempotencyKey string    `json:"idempotency_key,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// RecordUsageParams are the parameters for recording usage.
type RecordUsageParams struct {
	SubscriptionID string     `json:"subscription_id"`
	ProductKey     string     `json:"product_key,omitempty"`
	Quantity       int        `json:"quantity"`
	Timestamp      *time.Time `json:"timestamp,omitempty"`
	IdempotencyKey string     `json:"idempotency_key,omitempty"`
}

// RecordUsage records usage for a metered subscription.
func (s *BillingService) RecordUsage(ctx context.Context, params *RecordUsageParams) (*UsageRecord, error) {
	var record UsageRecord
	err := s.client.request(ctx, "POST", "/billing/v1/usage", params, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// RecordUsageByKey is a convenience method for recording usage with individual parameters.
func (s *BillingService) RecordUsageByKey(ctx context.Context, subscriptionID, productKey string, quantity int, idempotencyKey string) (*UsageRecord, error) {
	return s.RecordUsage(ctx, &RecordUsageParams{
		SubscriptionID: subscriptionID,
		ProductKey:     productKey,
		Quantity:       quantity,
		IdempotencyKey: idempotencyKey,
	})
}

// UsageSummary is an aggregated usage summary.
type UsageSummary struct {
	SubscriptionID string `json:"subscription_id"`
	ProductKey     string `json:"product_key"`
	TotalUsage     int    `json:"total_usage"`
	RecordCount    int    `json:"record_count"`
	PeriodStart    string `json:"period_start"`
	PeriodEnd      string `json:"period_end"`
}

// GetUsageSummaryParams are the parameters for getting a usage summary.
type GetUsageSummaryParams struct {
	SubscriptionID string
	ProductKey     string
}

// GetUsageSummary gets aggregated usage for a subscription.
func (s *BillingService) GetUsageSummary(ctx context.Context, params *GetUsageSummaryParams) (*UsageSummary, error) {
	path := "/billing/v1/usage?subscription_id=" + params.SubscriptionID
	if params.ProductKey != "" {
		path += "&product_key=" + params.ProductKey
	}

	var summary UsageSummary
	err := s.client.request(ctx, "GET", path, nil, &summary)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

// GetUsageSummaryByKey is a convenience method for getting usage with individual parameters.
func (s *BillingService) GetUsageSummaryByKey(ctx context.Context, subscriptionID, productKey string) (*UsageSummary, error) {
	return s.GetUsageSummary(ctx, &GetUsageSummaryParams{
		SubscriptionID: subscriptionID,
		ProductKey:     productKey,
	})
}

// ============================================================
// PORTAL
// ============================================================

// PortalLink represents a customer portal link.
type PortalLink struct {
	PortalURL string    `json:"portal_url"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// CreatePortalLinkParams are the parameters for creating a portal link.
type CreatePortalLinkParams struct {
	ExpiresInHours int `json:"expires_in_hours,omitempty"`
}

// CreatePortalLink creates a portal link for a customer.
func (s *BillingService) CreatePortalLink(ctx context.Context, customerID string, params *CreatePortalLinkParams) (*PortalLink, error) {
	var link PortalLink
	err := s.client.request(ctx, "POST", "/billing/v1/customers/"+customerID+"/portal-link", params, &link)
	if err != nil {
		return nil, err
	}
	return &link, nil
}
