package report

type ReportQuery struct {
	ReportName string
	DimensionName string
	DimensionValues []string
	Queries []NamedQuery
}

type NamedQuery struct {
	Name string
	Query string
	Params []QueryParam
	ResultType string
}

type SumType struct {
	Sum int
}

type CountType struct {
	Count int
}

type NameSumType struct {
	Name string
	Count int
}

type QueryParam struct {
	Name string
	Type string
}

type IntParam struct {
	Name string
	Value int
}

type StringParam struct {
	Name string
	Value string
}

type DailyReport struct {
	Name string
	Summaries []DailySummary
}

type WeeklyReport struct {
	Name string
	Summaries []WeeklySummary
}

type MonthlyReport struct {
	Name string
	Summaries []MonthlySummary
}

type DailySummary struct {
	Name string
	DayOfMonth int
	DayOfWeek string
	Year int
	Month int
	Count int
}

type WeeklySummary struct {
	Name string
	Category string
	WeekNumber int
	Year int
	Month int
	Count int
	DailySummaries []DailySummary
}


type MonthlySummary struct {
	Name string
	Category string
	Month string
	Year int
	Count int
	WeeklySummaries []WeeklySummary
}

type CreateReportRequest struct {
	Name            string       `form:"Name" json:"Name" binding:"required"`
	Query           string       `form:"Query" json:"Query" binding:"required"`
	Params          []QueryParam `form:"Params" json:"Params" binding:"required"`
	ResultType      string       `form:"ResultType" json:"ResultType" binding:"required"`
	DimensionName   string       `form:"DimensionName" json:"DimensionName" binding:"required"`
	DimensionValues []string     `form:"DimensionValues" json:"DimensionValues" binding:"required"`
}

type Query struct {
	Id              string
	Query           string
	QueryParams     []QueryParam
	DimensionName   string
	DimensionValues []string
}

