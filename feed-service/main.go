package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/lucasvillalbaar/go-events-cqrs/database"
	"github.com/lucasvillalbaar/go-events-cqrs/events"
	"github.com/lucasvillalbaar/go-events-cqrs/repository"
)

type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NAT_ADDRESS"`
}

func main() {
	var cfg Config

	err := envconfig.Process("", cfg)

	if err != nil {
		log.Fatalf("%v", err)
	}

	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)

	repo, err := database.NewPostgresRepository(addr)

	if err != nil {
		log.Fatalf("%v", err)
	}

	repository.SetRepository(repo)

	n, err := events.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))

	if err != nil {
		log.Fatalf("%v", err)
	}

	events.SetEventStore(n)

	defer events.Close()

	router := newRouter()

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("%v", err)
	}

}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()

	router.HandleFunc("/feed", createdFeedHandler).Methods(http.MethodPost)
	return router
}