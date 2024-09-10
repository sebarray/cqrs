package main

import (
	event "cqrs/events"
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var Config Config
	err := envconfig.Process("", &Config)
	if err != nil {
		log.Fatal(err.Error())
	}

	hub := NewHub()

	n, err := event.NewNatsEventStore(fmt.Sprintf("nats://%s", Config.NatsAddress))
	if err != nil {
		log.Fatal(err.Error())
	}

	err = n.OnCreatedFeed(func(m event.MessageFeedCreated) {
		hub.Broadcast(newCreatedFeedMessage(m.ID, m.Title, m.Description, m.CreatedAt), nil)
	})
	if err != nil {
		log.Fatal(err.Error())

	}
	event.SetEventStore(n)

	defer event.Close()
	go hub.Run()
	http.HandleFunc("/ws", hub.HandleWebSocket)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}

}
