package report

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

func Schedule(query Query, cron string) (*gocron.Job, error) {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.StartAsync()
	switch cron {
	case "MONTHLY":
		job, err := scheduler.Every(1).Month(1).Do(GenerateReport, query)
		if err != nil {
			log.Printf("Error add query to schedule. Query %s %s", query.Id, err.Error())
			return nil, err
		}
		return job, err
	case "WEEKLY":
		job, err := scheduler.Every(1).Week().Do(GenerateReport, query)
		if err != nil {
			log.Printf("Error add query to schedule. Query %s %s", query.Id, err.Error())
			return nil, err
		}
		return job, err
	case "DAILY":
		job, err := scheduler.Every(1).Day().Do(GenerateReport, query)
		if err != nil {
			log.Printf("Error add query to schedule. Query %s %s", query.Id, err.Error())
			return nil, err
		}
		return job, err
	case "HOURLY":
		job, err := scheduler.Every(1).Hour().Do(GenerateReport, query)
		if err != nil {
			log.Printf("Error add query to schedule. Query %s %s", query.Id, err.Error())
			return nil, err
		}
		return job, err
	case "MINUTE":
		job, err := scheduler.Every(1).Minute().Do(GenerateReport, query)
		if err != nil {
			log.Printf("Error add query to schedule. Query %s %s", query.Id, err.Error())
			return nil, err
		}
		return job, err
	default:
		return nil, fmt.Errorf("schedule %s not supported", cron)
	}
}
