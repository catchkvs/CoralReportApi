package report

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
   "cloud.google.com/go/bigquery"
)

var client *bigquery.Client

func init() {
	log.Println("Initialize report generator")
	ctx := context.Background()
	projectID := ""
	client, _ = bigquery.NewClient(ctx, projectID)
}

func GenerateReport(reportName string) {
	reportFile := "resource/" + reportName + ".json"
	data, err := ioutil.ReadFile(reportFile)
	if err != nil {
		log.Fatal(err)
	}
	reportQuery := ReportQuery{}
	json.Unmarshal(data, &reportQuery)
	log.Println(reportQuery)
}


