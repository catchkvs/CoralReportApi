package report

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"k8scale.io/coral/reportgen/pkg/config"
	"log"
	"strings"
)

type QueryStorage struct {
	queries []*ScheduledQuery
	firestore *firestore.Client
	ctx context.Context
}

func NewQueryStorage() *QueryStorage {
	projectId := config.GetProperty("coral.reportapi.projectid")
	ctx = context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		log.Printf("error initializing storage client")
	}
	storage := QueryStorage{
		queries:   make([]*ScheduledQuery, 0),
		firestore: client,
		ctx:       ctx,
	}
	return &storage
}

func (storage *QueryStorage) AddQuery(scheduledQuery *ScheduledQuery) error {
	_, err := storage.firestore.Collection("scheduled-query").Doc(scheduledQuery.Query.Id).Set(storage.ctx, scheduledQuery.Query)
	if err != nil {
		log.Printf("error storing scheduled query into storage %s. error - %s", scheduledQuery.Query.Id, err.Error())
	}
	storage.queries = append(storage.queries, scheduledQuery)
	return nil
}

func (storage *QueryStorage) GetAllStoredQueries() ([]Query, error) {
	results := storage.firestore.Collection("scheduled-query").Documents(storage.ctx)
	queries := make([]Query, 0)
	for {
		res, err := results.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("failed to iterate")
			return nil, err
		}
		var query Query

		err = res.DataTo(&query)
		if err != nil {
			log.Printf("Error deserializing stored data to query %s", err.Error())
			return nil, err
		}
		queries = append(queries, query)
	}
	log.Printf("Read %d queries from database", len(queries))
	return queries, nil
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

func (storage *QueryStorage) DeleteQuery(id string) error {
	_, err := storage.firestore.Collection("scheduled-query").Doc(id).Delete(storage.ctx)
	if err != nil {
		log.Printf("Error deleting scheduled query %s, error - %s", id, err.Error())
	}
	return err
}
