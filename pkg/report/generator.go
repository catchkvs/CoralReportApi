package report

import (
	"context"
	"encoding/json"
	"log"
   "cloud.google.com/go/bigquery"
)

func init() {
	log.Println("Creating storage bucket instance")
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		// TODO: Handle error.
	}
}


