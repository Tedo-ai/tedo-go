# Tedo Go SDK

Official Go client for the [Tedo API](https://tedo.ai/docs).

## Installation

```bash
go get github.com/tedo-ai/tedo-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/tedo-ai/tedo-go"
)

func main() {
    client := tedo.NewClient("tedo_live_xxx")

    // Create a customer
    customer, err := client.Billing.CreateCustomer(context.Background(), &tedo.CreateCustomerParams{
        Email: "user@example.com",
        Name:  "Acme Corp",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created customer: %s\n", customer.ID)

    // Create a subscription
    subscription, err := client.Billing.CreateSubscription(context.Background(), &tedo.CreateSubscriptionParams{
        CustomerID: customer.ID,
        PriceID:    "price_xxx",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created subscription: %s\n", subscription.ID)

    // Check entitlement
    result, err := client.Billing.CheckEntitlement(context.Background(), &tedo.CheckEntitlementParams{
        CustomerID:     customer.ID,
        EntitlementKey: "api_access",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Has access: %v\n", result.HasAccess)
}
```

## Configuration

### Custom Base URL

```go
client := tedo.NewClient("tedo_live_xxx").
    WithBaseURL("https://api.staging.tedo.ai/v1")
```

### Custom HTTP Client

```go
httpClient := &http.Client{
    Timeout: 60 * time.Second,
}

client := tedo.NewClient("tedo_live_xxx").
    WithHTTPClient(httpClient)
```

## Error Handling

```go
customer, err := client.Billing.GetCustomer(ctx, "cus_nonexistent")
if err != nil {
    if tedo.IsNotFound(err) {
        fmt.Println("Customer not found")
        return
    }
    if tedo.IsValidationError(err) {
        fmt.Printf("Validation error: %v\n", err)
        return
    }
    log.Fatal(err)
}
```

## Pagination

```go
var allCustomers []tedo.Customer
var cursor string

for {
    list, err := client.Billing.ListCustomers(ctx, &tedo.ListCustomersParams{
        Limit:  100,
        Cursor: cursor,
    })
    if err != nil {
        log.Fatal(err)
    }

    allCustomers = append(allCustomers, list.Customers...)

    if list.NextCursor == "" {
        break
    }
    cursor = list.NextCursor
}
```

## Available Services

### Billing

| Method | Description |
|--------|-------------|
| `CreateCustomer` | Create a new customer |
| `GetCustomer` | Get a customer by ID |
| `ListCustomers` | List all customers |
| `UpdateCustomer` | Update a customer |
| `DeleteCustomer` | Delete a customer |
| `CreateSubscription` | Create a subscription |
| `GetSubscription` | Get a subscription |
| `CancelSubscription` | Cancel a subscription |
| `CheckEntitlement` | Check feature access |
| `RecordUsage` | Record metered usage |
| `GetUsageSummary` | Get usage summary |

## License

MIT
