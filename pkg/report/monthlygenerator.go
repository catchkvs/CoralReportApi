package report

import (
	"context"
	"encoding/json"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"log"
   "cloud.google.com/go/bigquery"
	"strings"
)

var client *bigquery.Client
var ctx = context.Background()

func init() {
	log.Println("Initialize report generator")
	projectID := ""
	client, _ = bigquery.NewClient(ctx, projectID)
}

func GenerateCurrentMonthReport(reportName string, params map[string] string) {
	reportFile := "resource/" + reportName + ".json"
	data, err := ioutil.ReadFile(reportFile)
	if err != nil {
		log.Fatal(err)
	}
	reportQuery := ReportQuery{}
	json.Unmarshal(data, &reportQuery)
	log.Println(reportQuery)
	for idx, query := range reportQuery.Queries {
		log.Println("Running the query : ", idx)
		for _, param := range query.Params {
			paramStr := "${" + param.Name + "}"
			query.Query = strings.Replace(query.Query, paramStr, params[paramStr], 1 )
		}
		RunQuery(query.Query)
	}
}

func RunQuery(query string) {
	q := client.Query(query)
	q.Read(ctx)
	it, err := q.Read(ctx)
	if err != nil {
		log.Print(err)
	}
	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}
		log.Println(values)
	}
}


