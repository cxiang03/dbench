package dbench

import (
	"context"

	"github.com/meilisearch/meilisearch-go"
)

type Meili struct {
	*meilisearch.Client
}

func NewMeili() *Meili {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://103.3.60.74:7700",
		APIKey: "123qQ123",
	})
	return &Meili{Client: client}
}

func (m *Meili) UpdateFilterableAttributes(_ context.Context) error {
	attributes := []string{"price", "time_stamp", "postcode", "p_type", "is_new", "duration"}
	_, err := m.Index("prices").UpdateFilterableAttributes(&attributes)
	return err
}

func (m *Meili) Insert(_ context.Context, records []*Record) error {
	_, err := m.Index("prices").AddDocuments(records, "uuid")
	return err
}
