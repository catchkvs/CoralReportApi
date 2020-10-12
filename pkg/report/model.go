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