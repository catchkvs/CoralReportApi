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

func GenerateReport(query Query) {
	log.Printf("Initialize report generator for query %s", query)
	projectId := config.GetProperty("coral.reportapi.projectid")
	var err error
	client, err = bigquery.NewClient(ctx, projectId)
	if err != nil {
		log.Printf("Error initializing client %s", err.Error())
		return
	}
	replaceCurrentTimeBaseWildcards(&query)
	log.Printf("Query Dimension name %s type %s values %v", query.Dimension.Name, query.Dimension.Type, query.Dimension.Values)
	if &query.Dimension != nil {
		log.Printf("Query has multiple dimensions %s", query.Dimension.Name)
		switch query.Dimension.Type {
		case "STRING":
			for _, dv := range query.Dimension.Values {
				query.Query = strings.Replace(query.Query, "${"+query.Dimension.Name+"}", "'"+dv+"'", 1)
				log.Printf("Final query after all replacements %s", query.Query)
				runQuery(query, query.Dimension.Name, dv)
			}
		case "INT":
			for _, dv := range query.Dimension.Values {
				query.Query = strings.Replace(query.Query, "${"+query.Dimension.Name+"}", dv, 1)
				log.Printf("Final query after all replacements %s", query.Query)
				runQuery(query, query.Dimension.Name, dv)
			}
		default:
			log.Printf("Dimension %s not supported", query.Dimension.Type)
		}
	} else {
		log.Printf("No Dimension was found, final query %s", query.Query)
		runQuery(query, "", "")
	}
}

func replaceCurrentTimeBaseWildcards(query *Query) {
	WildCards := [7]string{"${CURRENT_YEAR}", "${CURRENT_MONTH}", "${CURRENT_DAY_OF_MONTH}", "${CURRENT_DAY_OF_WEEK}",
		"${CURRENT_HOUR_OF_DAY}", "${CURRENT_MINUTE_OF_HOUR}", "${CURRENT_SECOND_OF_MINUTE}"}
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
}

func runQuery(query Query, rangeColumn string, rangeKey string) {
	q := client.Query(query.Query)
	it, err := q.Read(ctx)
	if err != nil {
		log.Printf("Error running query %s - %s", query, err.Error())
		return
	}
	var result []map[string]bigquery.Value
	for {
		row := make(map[string]bigquery.Value)
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error reading iterator %s", err.Error())
			return
		}
		result = append(result, row)
	}
	content, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshalling records to json %s", err.Error())
		return
	}
	log.Printf("Query result %s", content)
	if rangeColumn != "" {
		err = Put(query.Id+"#"+rangeColumn+"#"+rangeKey, string(content))
	} else {
		err = Put(query.Id, string(content))
	}
	if err != nil {
		log.Printf("Error storing query result for id %s %s", query.Id, err.Error())
		return
	}
}
