package repository

import (
	"context"
	model "cqrs/models"
)

type Repository interface {
	Close()
	InsertFeed(ctx context.Context, feed model.Feed) error
	ListFeed(ctx context.Context) ([]model.Feed, error)
}

var repository Repository

func SetRepository(repo Repository) {
	repository = repo
}
func Close() {
	repository.Close()
}
func InsertFeed(ctx context.Context, feed model.Feed) error {
	return repository.InsertFeed(ctx, feed)
}

func ListFeed(ctx context.Context) ([]model.Feed, error) {
	return repository.ListFeed(ctx)
}
