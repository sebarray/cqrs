package event

import (
	"context"
	model "cqrs/models"
)

type EventStore interface {
	Close()
	Publish(ctx context.Context, feed *model.Feed) error
	Subscribe(ctx context.Context, feed *model.Feed) (<-chan MessageFeedCreated, error)
	OnCreatedFeed(f func(MessageFeedCreated)) error
}

var eventStore EventStore

func SetEventStore(es EventStore) {
	eventStore = es
}

func Close() {
	eventStore.Close()
}
func Publish(ctx context.Context, feed *model.Feed) error {
	return eventStore.Publish(ctx, feed)
}
func Subscribe(ctx context.Context, feed *model.Feed) (<-chan MessageFeedCreated, error) {
	return eventStore.Subscribe(ctx, feed)
}
func OnCreatedFeed(f func(MessageFeedCreated)) error {
	return eventStore.OnCreatedFeed(f)
}
