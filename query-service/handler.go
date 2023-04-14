package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/lucasvillalbaar/go-events-cqrs/events"
	"github.com/lucasvillalbaar/go-events-cqrs/models"
	"github.com/lucasvillalbaar/go-events-cqrs/repository"
	"github.com/lucasvillalbaar/go-events-cqrs/search"
)

func onCreatedFeed(m events.CreatedFeedMessage) {
	feed := models.Feed{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}

	if err := search.IndexFeed(context.Background(), &feed); err != nil {
		log.Print("failed to index feed:%v", err)
	}
}

func listFeedsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	feeds, err := repository.ListFeeds(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)

}

func searcHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	query := r.URL.Query().Get("q")

	if len(query) == 0 {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	feeds, err := search.SearchFeed(ctx, query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}
