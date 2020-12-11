package report

import (
	"github.com/google/uuid"
	"log"
	"strings"
)

type QueryStorage struct {
	queries []*Query
}

func NewQueryStorage() *QueryStorage {
	return &QueryStorage{
		queries: make([]*Query, 0),
	}
}

func (storage *QueryStorage) AddQuery(query string, params []QueryParam, dimensionName string, dimensionValues []string) Query {
	random, _ := uuid.NewRandom()
	randomId := random.String()
	q := Query{Id: randomId, Query: query, QueryParams: params, DimensionName: dimensionName, DimensionValues: dimensionValues}
	storage.queries = append(storage.queries, &q)
	return q
}

func (storage *QueryStorage) GetQuery(id string) *Query {
	log.Printf("queries present right now %d", len(storage.queries))
	for _, query := range storage.queries {
		log.Printf("Id - %s. Incoming Id %s", query.Id, id)
		if strings.Compare(id, query.Id) == 0 {
			return query
		}
	}
	log.Printf("No query found with Id %s", id)
	return &Query{}
}
