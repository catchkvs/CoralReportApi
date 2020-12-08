package report

import (
	"cloud.google.com/go/bigquery"
	"context"
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
	RunQuery(query.Query)
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
