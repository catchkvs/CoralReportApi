package report

import (
	"log"
	"strings"
)

type QueryStorage struct {
	queries []*ScheduledQuery
}

func NewQueryStorage() *QueryStorage {
	return &QueryStorage{
		queries: make([]*ScheduledQuery, 0),
	}
}

func (storage *QueryStorage) AddQuery(scheduledQuery *ScheduledQuery) {
	storage.queries = append(storage.queries, scheduledQuery)
}

func (storage *QueryStorage) GetQuery(id string) *ScheduledQuery {
	log.Printf("queries present right now %d", len(storage.queries))
	for _, scheduledQuery := range storage.queries {
		log.Printf("Id - %s. Incoming Id %s", scheduledQuery.Query.Id, id)
		if strings.Compare(id, scheduledQuery.Query.Id) == 0 {
			return scheduledQuery
		}
	}
	log.Printf("No scheduledQuery found with Id %s", id)
	return &ScheduledQuery{}
}
