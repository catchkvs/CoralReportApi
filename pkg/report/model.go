package report

import "github.com/go-co-op/gocron"

type NamedQuery struct {
	Name       string
	Query      string
	Params     []QueryParam
	ResultType string
}

type SumType struct {
	Sum int
}

type CountType struct {
	Count int
}

type NameSumType struct {
	Name  string
	Count int
}

type QueryParam struct {
	Name string
	Type string
}

type IntParam struct {
	Name  string
	Value int
}

type StringParam struct {
	Name  string
	Value string
}

type DailyReport struct {
	Name      string
	Summaries []DailySummary
}

type WeeklyReport struct {
	Name      string
	Summaries []WeeklySummary
}

type MonthlyReport struct {
	Name      string
	Summaries []MonthlySummary
}

type DailySummary struct {
	Name       string
	DayOfMonth int
	DayOfWeek  string
	Year       int
	Month      int
	Count      int
}

type WeeklySummary struct {
	Name           string
	Category       string
	WeekNumber     int
	Year           int
	Month          int
	Count          int
	DailySummaries []DailySummary
}

type MonthlySummary struct {
	Name            string
	Category        string
	Month           string
	Year            int
	Count           int
	WeeklySummaries []WeeklySummary
}

type CreateReportRequest struct {
	Name       string
	Query      string
	ResultType string
	Dimension  QueryDimension
	Cron       string
}

type QueryDimension struct {
	Name   string
	Type   string
	Values []string
}

type Query struct {
	Id        string
	Query     string
	Dimension QueryDimension
	Cron      string
}

type ScheduledQuery struct {
	Query        Query
	ScheduledJob *gocron.Job
}

type ViewReportRequest struct {
	Id             string
	DimensionName  string
	DimensionValue string
}

type GenerateReportRequest struct {
	Id string
}

type DeleteReportRequest struct {
	Id string
}
