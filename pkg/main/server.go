package main

import (
	"cloud.google.com/go/bigquery"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"k8scale.io/coral/reportgen/pkg/report"
	"log"
	"net/http"
)

func main() {
	storage := report.NewQueryStorage()
	router := gin.Default()
	router.POST("/api/v1/create-report", func(context *gin.Context) {
		var request report.CreateReportRequest
		if err := context.ShouldBindJSON(&request); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		query := storage.AddQuery(request.Query, request.Params, request.DimensionName, request.DimensionValues)
		log.Printf("Stored query %s", query.Id)
		scheduledQuery := report.ScheduledQuery{Query: query, Cron: request.Cron}
		if err := report.Schedule(scheduledQuery); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		context.JSON(http.StatusOK, gin.H{"Id": query.Id})
	})

	router.GET("/api/v1/generate-report", func(context *gin.Context) {
		var request ReportGenerateRequest
		if err := context.ShouldBindQuery(&request); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		query := storage.GetQuery(request.Id)
		report.GenerateCurrentMonthReport(*query)
	})

	router.GET("/api/v1/view-report", func(context *gin.Context) {
		var request ViewQueryReportRequest
		if err := context.ShouldBindQuery(&request); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		result, err := report.Get(request.Id)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		var decoded []map[string]bigquery.Value
		err = json.Unmarshal([]byte(result), &decoded)
		if err != nil {
			log.Printf("Error unmarshalling the stored result")
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		context.JSON(http.StatusOK, gin.H{"report": decoded})
	})
	log.Fatal(router.Run(":8080"))
}

type ReportGenerateRequest struct {
	Id string `form:"Id" json:"Id" binding:"required"`
}

type ViewQueryReportRequest struct {
	Id string `form:"Id" json:"Id" binding:"required"`
}