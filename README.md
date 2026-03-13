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

    // Billing: create a customer
    customer, err := client.Billing.CreateCustomer(context.Background(), &tedo.CreateCustomerParams{
        Email: "user@example.com",
        Name:  "Acme Corp",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created customer: %s\n", customer.ID)

    // Sales: create a pipeline and a lead
    pipeline, err := client.Sales.CreatePipeline(context.Background(), &tedo.CreatePipelineParams{
        Name:         "Inbound",
        ResourceType: tedo.ResourceTypeLead,
    })
    if err != nil {
        log.Fatal(err)
    }

    lead, err := client.Sales.CreateLead(context.Background(), &tedo.CreateLeadParams{
        Label:      "Acme Corp",
        PipelineID: pipeline.ID,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created lead: %s\n", lead.ID)
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

## Services

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

### Sales

The Sales service covers pipelines, stages, leads, deals, activities, notes, and contacts.

#### Constants

**Activity types** (`ActivityType*`):

| Constant | Value |
|----------|-------|
| `ActivityTypeTask` | `"task"` |
| `ActivityTypeCall` | `"call"` |
| `ActivityTypeEmail` | `"email"` |
| `ActivityTypeMeeting` | `"meeting"` |
| `ActivityTypeDeadline` | `"deadline"` |

**Pipeline resource types** (`ResourceType*`):

| Constant | Value |
|----------|-------|
| `ResourceTypeLead` | `"lead"` |
| `ResourceTypeDeal` | `"deal"` |

**Stage outcomes** (`Outcome*`):

| Constant | Value |
|----------|-------|
| `OutcomePositive` | `"positive"` |
| `OutcomeNegative` | `"negative"` |

#### Methods

| Method | Description |
|--------|-------------|
| `CreatePipeline` | Create a pipeline for leads or deals |
| `GetPipeline` | Get a pipeline by ID |
| `ListPipelines` | List all pipelines |
| `UpdatePipeline` | Update a pipeline |
| `DeletePipeline` | Delete a pipeline |
| `CreateStage` | Create a stage in a pipeline |
| `GetStage` | Get a stage by ID |
| `ListStages` | List stages for a pipeline |
| `UpdateStage` | Update a stage |
| `DeleteStage` | Delete a stage |
| `CreateLead` | Create a new lead |
| `GetLead` | Get a lead by ID |
| `ListLeads` | List leads, optionally filtered by pipeline |
| `UpdateLead` | Update a lead |
| `DeleteLead` | Delete a lead |
| `MoveLeadStage` | Move a lead to a different stage |
| `ConvertLeadToDeal` | Convert a lead into a deal |
| `CreateDeal` | Create a new deal |
| `GetDeal` | Get a deal by ID |
| `ListDeals` | List deals, optionally filtered by pipeline |
| `UpdateDeal` | Update a deal |
| `DeleteDeal` | Delete a deal |
| `MoveDealStage` | Move a deal to a different stage |
| `CreateActivity` | Create an activity linked to entities |
| `GetActivity` | Get an activity by ID |
| `ListActivities` | List activities, optionally filtered by type or completion |
| `UpdateActivity` | Update an activity |
| `DeleteActivity` | Delete an activity |
| `CompleteActivity` | Mark an activity as completed or uncompleted |
| `CreateSalesNote` | Create a note linked to sales entities |
| `GetSalesNote` | Get a note by ID |
| `ListSalesNotes` | List all notes |
| `UpdateSalesNote` | Update a note |
| `DeleteSalesNote` | Delete a note |
| `CreateContactBase` | Create a new contact base |
| `GetContactBase` | Get a contact base by ID |
| `ListContactBases` | List all contact bases |
| `CreatePerson` | Create a person in a contact base |
| `GetPerson` | Get a person by ID |
| `ListPersons` | List persons in a contact base |
| `UpdatePerson` | Update a person |
| `DeletePerson` | Delete a person |
| `CreateOrganization` | Create an organization in a contact base |
| `GetOrganization` | Get an organization by ID |
| `ListOrganizations` | List organizations in a contact base |
| `UpdateOrganization` | Update an organization |
| `DeleteOrganization` | Delete an organization |

#### Examples

**Create a pipeline with stages:**

```go
pipeline, err := client.Sales.CreatePipeline(ctx, &tedo.CreatePipelineParams{
    Name:         "Sales",
    ResourceType: tedo.ResourceTypeDeal,
})

outcome := tedo.OutcomePositive
client.Sales.CreateStage(ctx, pipeline.ID, &tedo.CreateStageParams{
    Name:       "Closed Won",
    Position:   3,
    IsTerminal: true,
    Outcome:    &outcome,
})
```

**Manage leads:**

```go
lead, err := client.Sales.CreateLead(ctx, &tedo.CreateLeadParams{
    Label:      "Acme Corp",
    PipelineID: pipeline.ID,
    StageID:    stage.ID,
})

// Move to next stage
lead, err = client.Sales.MoveLeadStage(ctx, lead.ID, nextStage.ID)

// Convert to deal
deal, err := client.Sales.ConvertLeadToDeal(ctx, lead.ID, &tedo.ConvertLeadParams{
    DealPipelineID: dealPipeline.ID,
    DealStageID:    dealStage.ID,
    DealLabel:      "Acme Corp - Enterprise",
})
```

**Create activities with links:**

Use the link helpers to attach an activity to one or more sales entities:

```go
activity, err := client.Sales.CreateActivity(ctx, &tedo.CreateActivityParams{
    Type:    tedo.ActivityTypeCall,
    Subject: "Discovery call",
    Links: []tedo.ActivityLink{
        tedo.LeadLink(lead.ID),           // primary link
        tedo.PersonLink(person.ID),        // secondary link
        tedo.OrganizationLink(org.ID),     // secondary link
    },
})

// Mark as done
client.Sales.CompleteActivity(ctx, activity.ID, true)
```

Available link helpers: `LeadLink`, `DealLink`, `PersonLink`, `OrganizationLink`.
For notes: `LeadNoteLink`, `DealNoteLink`, `PersonNoteLink`, `OrganizationNoteLink`.

**Contacts:**

```go
// List contact bases (each workspace has a default one)
bases, _ := client.Sales.ListContactBases(ctx)
baseID := bases.ContactBases[0].ID

person, _ := client.Sales.CreatePerson(ctx, baseID, &tedo.CreatePersonParams{
    FirstName: "Jane",
    LastName:  "Smith",
    Email:     "jane@acme.com",
})

org, _ := client.Sales.CreateOrganization(ctx, base.ID, &tedo.CreateOrganizationParams{
    Name:    "Acme Corp",
    Website: ptr("https://acme.com"),
})
```

## License

MIT
