package main

import (
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
			queryId := storage.AddQuery(request.Query, request.Params, request.DimensionName, request.DimensionValues)
			log.Printf("Stored query %s", queryId)
			context.JSON(http.StatusOK, gin.H{"Id": queryId})
	})
	router.GET("/api/v1/generate-report", func(context *gin.Context) {
		var request ReportGenerateRequest
		if err := context.ShouldBindQuery(&request); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		query := storage.GetQuery(request.Id)
		report.GenerateCurrentMonthReport(*query)
	})
	router.GET("/api/v1/view-report", viewReportHandler)
	log.Fatal(router.Run(":8080"))
}

type ReportGenerateRequest struct {
	Id string `form:"Id" json:"Id" binding:"required"`
}

func reportCreateHandler(context *gin.Context) {

}





func viewReportHandler(context *gin.Context) {
}


