package search

import (
	"bytes"
	"context"
	model "cqrs/models"
	"encoding/json"

	elastic "github.com/elastic/go-elasticsearch/v7"
)

type ElasticSearchRepository struct {
	Client *elastic.Client
}

func NewElasticSearchRepository(url string) (*ElasticSearchRepository, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}

	return &ElasticSearchRepository{
		Client: client,
	}, nil
}

func (e *ElasticSearchRepository) Close() {
	//
}

func (e *ElasticSearchRepository) IndexFeed(ctx context.Context, feed model.Feed) error {
	body, _ := json.Marshal(feed)
	_, err := e.Client.Index("feeds", bytes.NewReader(body),
		e.Client.Index.WithDocumentID(feed.ID),
		e.Client.Index.WithContext(ctx),
		e.Client.Index.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	return nil
}

func (e *ElasticSearchRepository) SearchFeed(ctx context.Context, query string) ([]model.Feed, error) {
	var buf bytes.Buffer

	serchquery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":             query,
				"fields":            []string{"title", "description"},
				"fuzziness":         3,
				"cuttoff_frequency": 0.001,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(serchquery); err != nil {
		return nil, err
	}
	res, err := e.Client.Search(
		e.Client.Search.WithContext(ctx),
		e.Client.Search.WithIndex("feeds"),
		e.Client.Search.WithBody(&buf),
		e.Client.Search.WithTrackTotalHits(true),
		e.Client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, err
	}
	var eRes map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&eRes); err != nil {
		return nil, err
	}
	var feeds []model.Feed
	for _, hit := range eRes["hits"].(map[string]interface{})["hits"].([]interface{}) {

		feed := model.Feed{}
		sourse := hit.(map[string]interface{})["_source"]
		marshal, err := json.Marshal(sourse)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(marshal, &feed)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}
	return feeds, nil
}
