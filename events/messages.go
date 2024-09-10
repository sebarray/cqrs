package event

import "time"

type Message interface {
	Type() string
}

type MessageFeedCreated struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"create"`
}

func (m *MessageFeedCreated) Type() string {
	return "feed_created"
}
