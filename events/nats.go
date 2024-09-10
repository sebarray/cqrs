package event

import (
	"bytes"
	"context"
	model "cqrs/models"
	"encoding/gob"

	"github.com/nats-io/nats.go"
)

type NatsEventStore struct {
	Conn            *nats.Conn
	FeedCreatedSub  *nats.Subscription
	FeedCreatedChan chan MessageFeedCreated
}

func NewNatsEventStore(url string) (*NatsEventStore, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NatsEventStore{
		Conn: conn,
	}, nil
}

func (n *NatsEventStore) Close() {
	if n.Conn != nil {
		n.Conn.Close()
	}
	if n.FeedCreatedSub != nil {
		n.FeedCreatedSub.Unsubscribe()
	}

	close(n.FeedCreatedChan)
}
func (n *NatsEventStore) encodeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil

}

func (n *NatsEventStore) Publish(ctx context.Context, feed *model.Feed) error {
	msg := &MessageFeedCreated{
		ID:          feed.ID,
		Title:       feed.Title,
		Description: feed.Description,
		CreatedAt:   feed.CreatedAt,
	}
	b, err := n.encodeMessage(msg)
	if err != nil {
		return err
	}
	return n.Conn.Publish(msg.Type(), b)

}

func (n *NatsEventStore) dencodeMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	return gob.NewDecoder(&b).Decode(m)
}

func (n *NatsEventStore) OnCreatedFeed(f func(MessageFeedCreated)) (err error) {
	msg := MessageFeedCreated{}

	n.FeedCreatedSub, err = n.Conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		n.dencodeMessage(m.Data, &msg)
		f(msg)
	})
	return
}

func (n *NatsEventStore) Subscribe(ctx context.Context, feed *model.Feed) (<-chan MessageFeedCreated, error) {
	m := MessageFeedCreated{}
	n.FeedCreatedChan = make(chan MessageFeedCreated, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	n.FeedCreatedSub, err = n.Conn.ChanSubscribe(m.Type(), ch)
	if err != nil {
		return nil, err
	}
	go func() {

		for {
			select {
			case msg := <-ch:
				n.dencodeMessage(msg.Data, &m)
				n.FeedCreatedChan <- m

			}

		}

	}()

	return (<-chan MessageFeedCreated)(n.FeedCreatedChan), nil
}
