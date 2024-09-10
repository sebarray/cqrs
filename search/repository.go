package search

import (
	"context"
	model "cqrs/models"
)

type SearchRepository interface {
	Close()
	IndexFeed(ctx context.Context, feed model.Feed) error
	SearchFeed(ctx context.Context, query string) ([]model.Feed, error)
}

var searchRepo SearchRepository

func SetSearchRepository(repo SearchRepository) {
	searchRepo = repo
}

func Close() {
	searchRepo.Close()
}

func IndexFeed(ctx context.Context, feed model.Feed) error {
	return searchRepo.IndexFeed(ctx, feed)
}

func SearchFeed(ctx context.Context, query string) ([]model.Feed, error) {
	return searchRepo.SearchFeed(ctx, query)
}
