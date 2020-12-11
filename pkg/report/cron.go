package report

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

func Schedule(scheduledQuery ScheduledQuery) error {
	cron := scheduledQuery.Cron
	s1 := gocron.NewScheduler(time.UTC)
	s1.StartAsync()
	switch cron {
	case "MONTHLY":
		_, err := s1.Every(1).Month(1).Do(GenerateCurrentMonthReport, scheduledQuery.Query)
		if err != nil {
			log.Printf("Error add query to schedule. Query %s %s", scheduledQuery.Query.Id, err.Error())
		}
	case "WEEKLY":
		_, err := s1.Every(1).Week().Do(GenerateCurrentMonthReport, scheduledQuery.Query)
		if err != nil {
			log.Printf("Error add query to schedule. Query %s %s", scheduledQuery.Query.Id, err.Error())
			return err
		}
	case "DAILY":
		_, err := s1.Every(1).Day().Do(GenerateCurrentMonthReport, scheduledQuery.Query)
		if err != nil {
			log.Printf("Error add query to schedule. Query %s %s", scheduledQuery.Query.Id, err.Error())
			return err
		}
	case "HOURLY":
		_, err := s1.Every(1).Hour().Do(GenerateCurrentMonthReport, scheduledQuery.Query)
		if err != nil {
			log.Printf("Error add query to schedule. Query %s %s", scheduledQuery.Query.Id, err.Error())
			return err
		}
	case "MINUTE":
		_, err := s1.Every(1).Minute().Do(GenerateCurrentMonthReport, scheduledQuery.Query)
		if err != nil {
			log.Printf("Error add query to schedule. Query %s %s", scheduledQuery.Query.Id, err.Error())
			return err
		}
	default:
		return fmt.Errorf("Schedule %s not supported", cron)
	}
	return nil
}
