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

    // Projects: create a project and a work item
    project, err := client.Projects.CreateProject(context.Background(), &tedo.CreateProjectParams{
        Name:        "Q2 launch",
        Description: "Public API rollout",
    })
    if err != nil {
        log.Fatal(err)
    }

    workItem, err := client.Projects.CreateProjectWorkItem(context.Background(), project.ID, &tedo.CreateWorkItemParams{
        Title:       "Publish Go SDK bindings",
        Description: "Use the Projects public API",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created work item: %s\n", workItem.ID)

    // Tables: create a table and query editorial rows
    created, err := client.Tables.CreateTable(context.Background(), &tedo.CreateTableParams{
        Name: "Editorial Pipeline",
    })
    if err != nil {
        log.Fatal(err)
    }
    _, err = client.Tables.UpsertColumn(context.Background(), created.Table.ID, &tedo.UpsertColumnParams{
        Name: "Title",
        Key:  "title",
        Type: tedo.ColumnTypeText,
    })
    if err != nil {
        log.Fatal(err)
    }
    _, err = client.Tables.BulkUpsertRows(context.Background(), created.Table.ID, &tedo.BulkUpsertRowsParams{
        KeyColumn: "title",
        Rows: []tedo.JSONMap{{"title": "Customer story"}},
    })
    if err != nil {
        log.Fatal(err)
    }
    rows, err := client.Tables.QueryRows(context.Background(), created.Table.ID, &tedo.QueryRowsParams{
        Limit: 10,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found rows: %d\n", len(rows.Rows))
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
    Transport: &http.Transport{
        MaxIdleConnsPerHost: 10,
    },
}

client := tedo.NewClient("tedo_live_xxx").
    WithHTTPClient(httpClient)
```

For storage uploads and downloads, prefer request-scoped timeouts on the `context.Context` you pass into each call instead of setting `http.Client.Timeout`. That keeps long-running uploads cancellable without forcing a fixed global timeout.

### Retry Configuration

```go
client := tedo.NewClient("tedo_live_xxx").
    WithRetryConfig(tedo.RetryConfig{
        MaxRetries:     5,
        InitialBackoff: 250 * time.Millisecond,
        MaxBackoff:     3 * time.Second,
    })
```

Storage operations automatically retry transient transport failures and `408`, `429`, `500`, `502`, `503`, and `504` responses.

### Per-request Headers

Projects commands accept request options for idempotency and request tracing:

```go
project, err := client.Projects.CreateProject(ctx, &tedo.CreateProjectParams{
    Name: "Q2 launch",
}, tedo.WithIdempotencyKey("project-q2-launch"), tedo.WithRequestID("req_123"))
```

For Projects creates, deletes, and attachment attach/detach calls, the SDK generates an `Idempotency-Key` when you do not provide one.

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

### Projects

The Projects service covers projects, work items, workflow configuration, read-only comments, activity feeds, and file-reference attachments.

Projects V1 intentionally does not expose comment writes, raw multipart attachment upload, live collaborative edit fields, bulk operations, or task-named aliases. Attach files by uploading through Files/Storage first, then call `AttachFile` with the file ID.

#### Methods

| Method | Description |
|--------|-------------|
| `CreateProject` | Create a project |
| `GetProject` | Get a project by ID |
| `ListProjects` | List projects with cursor pagination |
| `UpdateProject` | Update a project |
| `ArchiveProject` | Archive a project |
| `RestoreProject` | Restore a project |
| `DeleteProject` | Delete a project |
| `CreateWorkItem` | Create a work item |
| `CreateProjectWorkItem` | Create a work item bound to a project path |
| `GetWorkItem` | Get a work item by ID |
| `ListWorkItems` | List work items with cursor pagination |
| `ListProjectWorkItems` | List work items in a project |
| `UpdateWorkItem` | Update a work item |
| `CompleteWorkItem` | Complete or reopen a work item |
| `ArchiveWorkItem` | Archive a work item |
| `RestoreWorkItem` | Restore a work item |
| `DeleteWorkItem` | Delete a work item |
| `ListSubtasks` | List work-item subtasks |
| `ListWorkItemActivity` | List work-item activity entries |
| `PeekNextDisplayID` | Preview the next display ID |
| `ListStatuses` | List workflow statuses |
| `CreateStatus` | Create a workflow status |
| `UpdateStatus` | Update a workflow status |
| `DeleteStatus` | Delete a workflow status |
| `ListWorkItemTypes` | List work-item types |
| `CreateWorkItemType` | Create a work-item type |
| `UpdateWorkItemType` | Update a work-item type |
| `DeleteWorkItemType` | Delete a work-item type |
| `ListPriorityLevels` | List priority levels |
| `UpdatePriorityLevel` | Update a priority level |
| `ResetPriorityLevel` | Reset a priority level |
| `ListComments` | List comments; V1 read-only |
| `ListAttachments` | List file-reference attachments |
| `AttachFile` | Attach an existing Files/Storage file reference |
| `DetachAttachment` | Detach a file reference |

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

### Storage

The Storage service covers buckets, objects, upload integrity checks, pre-signed URLs, and usage reporting.

| Method | Description |
|--------|-------------|
| `ListBuckets` | List buckets in the workspace |
| `CreateBucket` | Create a storage bucket |
| `GetBucket` | Get bucket metadata |
| `DeleteBucket` | Delete an empty bucket |
| `ListObjects` | List objects in a bucket |
| `PutObject` | Upload an object |
| `PutObjectWithOptions` | Upload an object with integrity options |
| `HeadObject` | Read object metadata without downloading the body |
| `GetObject` | Download an object |
| `DeleteObject` | Delete an object |
| `PresignURL` | Create a temporary download URL |
| `GetUsage` | Read storage usage totals |

**Upload with a content hash:**

```go
payload := []byte("hello world")
sum := sha256.Sum256(payload)

obj, err := client.Storage.PutObjectWithOptions(
    ctx,
    "bucket_123",
    "inbox/2026/04/message.eml",
    bytes.NewReader(payload),
    &tedo.PutObjectOptions{
        ContentType:   "message/rfc822",
        ContentSHA256: hex.EncodeToString(sum[:]),
    },
)
if err != nil {
    log.Fatal(err)
}

fmt.Println(obj.Hash)
```

**Check if an object already exists before uploading:**

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()

obj, err := client.Storage.HeadObject(ctx, "bucket_123", "sha256/abc123")
switch {
case err == nil:
    fmt.Printf("Already stored: %s (%d bytes)\n", obj.Key, obj.Size)
case tedo.IsNotFound(err):
    fmt.Println("Object not found, safe to upload")
default:
    log.Fatal(err)
}
```

## License

MIT
