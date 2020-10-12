package main

import (
	"k8scale.io/coral/reportgen/pkg/report"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/api/v1/generate-report", report.HandleReportGenerationRequest)
	http.HandleFunc("/api/v1/view-report", report.HandleReportGenerationRequest)
	log.Fatal(http.ListenAndServe(":3080", nil))
}
