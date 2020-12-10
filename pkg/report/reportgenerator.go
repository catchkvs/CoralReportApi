package report

import (
	"cloud.google.com/go/bigquery"
	"context"
	"encoding/json"
	"google.golang.org/api/iterator"
	"k8scale.io/coral/reportgen/pkg/config"
	"log"
	"strconv"
	"strings"
	"time"
)

var client *bigquery.Client
var ctx = context.Background()

func init() {
	log.Println("Initialize report generator")
	projectId := config.GetProperty("coral.reportapi.projectid")
	client, _ = bigquery.NewClient(ctx, projectId)
}

func GenerateCurrentMonthReport(query Query) {
	log.Println("Initialize report generator")
	projectId := config.GetProperty("coral.reportapi.projectid")
	var err error
	client, err = bigquery.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("Error initializing client %s", err.Error())
	}

	WildCards := [7]string{"${CURRENT_YEAR}", "${CURRENT_MONTH}", "${CURRENT_DAY_OF_MONTH}", "${CURRENT_DAY_OF_WEEK}", "${CURRENT_HOUR_OF_DAY}", "${CURRENT_MINUTE_OF_HOUR}", "${CURRENT_SECOND_OF_MINUTE}"}
	time.Now().Second()
	for _, wildcard := range WildCards {
		switch wildcard {
		case "${CURRENT_YEAR}":
			query.Query = strings.Replace(query.Query, wildcard, strconv.Itoa(time.Now().Year()), 1)
		case "${CURRENT_MONTH}":
			query.Query = strings.Replace(query.Query, wildcard, strconv.Itoa(int(time.Now().Month())), 1)
		case "${CURRENT_DAY_OF_MONTH}":
			query.Query = strings.Replace(query.Query, wildcard, strconv.Itoa(time.Now().Day()), 1)
		case "${CURRENT_HOUR_OF_DAY}":
			query.Query = strings.Replace(query.Query, wildcard, strconv.Itoa(time.Now().Hour()), 1)
		case "${CURRENT_MINUTE_OF_HOUR}":
			query.Query = strings.Replace(query.Query, wildcard, strconv.Itoa(time.Now().Minute()), 1)
		case "${CURRENT_SECOND_OF_MINUTE}":
			query.Query = strings.Replace(query.Query, wildcard, strconv.Itoa(time.Now().Second()), 1)
		}
	}
	log.Printf("Query after replacements %s", query.Query)
	RunQuery(query)
}

func RunQuery(query Query) {
	q := client.Query(query.Query)
	it, err := q.Read(ctx)
	if err != nil {
		log.Fatalf("Error running query %s - %s", query, err.Error())
	}
	var result []map[string]bigquery.Value
	for {
		row := make(map[string]bigquery.Value)
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error reading iterator %s", err.Error())
		}
		log.Printf("row %v", row)
		result = append(result, row)
	}
	content, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Error marshalling records to json %s", err.Error())
	}
	Put(query.Id, string(content))
}
