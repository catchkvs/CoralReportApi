package main

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"k8scale.io/coral/reportgen/pkg/report"
	"log"
	"net/http"
)

func main() {
	storage := report.NewQueryStorage()
	queries, err := storage.GetAllStoredQueries()
	scheduler := report.InitScheduler()
	if err != nil {
		log.Fatalf("error reading stored quries")
	}
	for _, query := range queries {
		job, err := report.Schedule(query)
		if err != nil {
			log.Fatalf("Error scheduling query %s", query.Id)
		}
		scheduledQuery := report.ScheduledQuery{Query: query, ScheduledJob: job}
		err = storage.AddQuery(&scheduledQuery)
		if err != nil {
			log.Printf("Error adding query %s", query.Id)
			log.Fatalf("Error scheduling query")
		}
	}

	log.Printf("Finished schuling already stored queries. Starting server")
	router := gin.Default()

	router.POST("/api/v1/create-report", func(context *gin.Context) {
		var request report.CreateReportRequest
		if err := context.ShouldBindJSON(&request); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "error deserializing request"})
			return
		}
		queryId := uuid.New().String()
		query := report.Query{Id: queryId, Query: request.Query, Dimension: request.Dimension, Cron: request.Cron}
		log.Printf("Creating query report %s", query.Query)
		job, err := report.Schedule(query)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "scheduling query"})
			return
		}
		scheduledQuery := report.ScheduledQuery{Query: query, ScheduledJob: job}
		err = storage.AddQuery(&scheduledQuery)
		if err != nil {
			log.Printf("Error adding query %s", query.Id)
			context.JSON(http.StatusInternalServerError, gin.H{"error": "scheduling query"})
		}
		log.Printf("created and scheduled query with id %s", queryId)
		context.JSON(http.StatusOK, gin.H{"Id": scheduledQuery.Query.Id})
	})

	router.GET("/api/v1/generate-report", func(context *gin.Context) {
		var request report.GenerateReportRequest
		if err := context.ShouldBindQuery(&request); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "deserializing request"})
			return
		}
		scheduledQuery := storage.GetQuery(request.Id)
		report.GenerateReport(scheduledQuery.Query)
	})

	router.GET("/api/v1/view-report", func(context *gin.Context) {
		var request report.ViewReportRequest
		if err := context.ShouldBindQuery(&request); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "deserializing request"})
			return
		}
		var result string
		var err error
		if request.DimensionName != "" && request.DimensionValue != "" {
			result, err = report.Get(request.Id + "#" + request.DimensionName + "#" + request.DimensionValue)
		} else {
			result, err = report.Get(request.Id)
		}
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "getting report"})
			return
		}
		if result == "" {
			context.JSON(http.StatusNotFound, gin.H{"error": "no report found with query id " + request.Id})
			return
		}
		var decoded []map[string]bigquery.Value
		err = json.Unmarshal([]byte(result), &decoded)
		if err != nil {
			log.Printf("Error unmarshalling the stored result %s", err.Error())
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing request"})
			return
		}
		context.JSON(http.StatusOK, gin.H{"report": decoded})
	})

	router.DELETE("/api/v1/delete-report", func(context *gin.Context) {
		var request report.DeleteReportRequest
		if err := context.ShouldBindQuery(&request); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "deserializing request"})
			return
		}
		scheduledQuery := storage.GetQuery(request.Id)
		if scheduledQuery == nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "no scheduled query found with id " + request.Id})
			return
		}
		err = storage.DeleteQuery(scheduledQuery.Query.Id)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting scheduled query"})
			return
		}
		scheduler.RemoveByReference(scheduledQuery.ScheduledJob)
		log.Printf("Removed scheduled query %s", scheduledQuery.Query.Id)
		context.JSON(http.StatusOK, gin.H{})
	})
	log.Fatal(router.Run(":8080"))
}

