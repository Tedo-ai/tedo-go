package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tedo "github.com/tedo-ai/tedo-go"
)

func main() {
	apiKey := os.Getenv("TEDO_API_KEY")
	if apiKey == "" {
		log.Fatal("TEDO_API_KEY is required")
	}

	ctx := context.Background()
	client := tedo.NewClient(apiKey)

	created, err := client.Tables.CreateTable(ctx, &tedo.CreateTableParams{
		Name: "Editorial Pipeline",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Tables.UpsertColumn(ctx, created.Table.ID, &tedo.UpsertColumnParams{
		Name: "Title",
		Key:  "title",
		Type: tedo.ColumnTypeText,
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Tables.UpsertColumn(ctx, created.Table.ID, &tedo.UpsertColumnParams{
		Name:       "Status",
		Key:        "status",
		Type:       tedo.ColumnTypeSelect,
		TypeConfig: tedo.JSONMap{"options": []any{"draft", "review", "published"}},
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Tables.BulkUpsertRows(ctx, created.Table.ID, &tedo.BulkUpsertRowsParams{
		KeyColumn: "title",
		Rows: []tedo.JSONMap{
			{"title": "Spring launch", "status": "draft"},
			{"title": "Customer story", "status": "review"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	rows, err := client.Tables.QueryRows(ctx, created.Table.ID, &tedo.QueryRowsParams{
		Filter: `eq(status,"review")`,
		Limit:  10,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range rows.Rows {
		fmt.Printf("%s %#v\n", row.ID, row.Data)
	}
}
