package search

import (
	"bytes"
	"context"
	"encoding/json"

	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/lucasvillalbaar/go-events-cqrs/models"
)

type ElasticSearchRepository struct {
	client *elastic.Client
}

func NewElastic(url string) (*ElasticSearchRepository, error) {
	client, err := elastic.NewClient(
		elastic.Config{
			Addresses: []string{url},
		})
	if err != nil {
		return nil, err
	}

	return &ElasticSearchRepository{client: client}, nil
}

func (r *ElasticSearchRepository) Close() {
	//No way to close an Elastic Client
}

func (r *ElasticSearchRepository) IndexFeed(ctx context.Context, feed *models.Feed) error {
	body, _ := json.Marshal(feed)
	_, err := r.client.Index(
		"feeds",
		bytes.NewReader(body),
		r.client.Index.WithDocumentID(feed.ID),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithRefresh("wait_for"),
	)

	return err
}
