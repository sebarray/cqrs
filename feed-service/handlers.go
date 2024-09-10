package main

import (
	event "cqrs/events"
	model "cqrs/models"
	"cqrs/repository"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"
)

type OnCreatedFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func OnCreatedFeedHandler(w http.ResponseWriter, r *http.Request) {
	var req OnCreatedFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	createAT := time.Now().UTC()
	id, err := ksuid.NewRandom()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	feed := model.Feed{
		ID:          id.String(),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   createAT,
	}
	if err := repository.InsertFeed(r.Context(), feed); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := event.Publish(r.Context(), &feed); err != nil {
		log.Println(err.Error())
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
	return

}
