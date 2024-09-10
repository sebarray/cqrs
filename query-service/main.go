package main

import (
	"cqrs/database"
	event "cqrs/events"
	"cqrs/repository"
	"cqrs/search"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDB          string `envconfig:"POSTGRES_DB"`
	PostgresUser        string `envconfig:"POSTGRES_USER"`
	PostgresPassword    string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress         string `envconfig:"NATS_ADDRESS"`
	ElasticSerchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/feeds", ListFeedHandler)
	router.HandleFunc("/search", SearchFeedHandler)
	return router
}

func main() {
	var Config Config
	err := envconfig.Process("", &Config)
	if err != nil {
		log.Fatal(err.Error())
	}
	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", Config.PostgresUser, Config.PostgresPassword, Config.PostgresDB)
	repo, err := database.NewPgRepository(addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	repository.SetRepository(repo)
	es, err := search.NewElasticSearchRepository(fmt.Sprintf("http://%s", Config.ElasticSerchAddress))

	if err != nil {
		log.Fatal(err.Error())
	}
	search.SetSearchRepository(es)
	defer search.Close()
	n, err := event.NewNatsEventStore(fmt.Sprintf("nats://%s", Config.NatsAddress))
	if err != nil {
		log.Fatal(err.Error())
	}

	err = n.OnCreatedFeed(OnCreatedFeed)
	if err != nil {
		log.Fatal(err.Error())

	}

	event.SetEventStore(n)
	defer event.Close()
	router := NewRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err.Error())
	}
}
