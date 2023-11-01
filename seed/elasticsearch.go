package dbench

import (
	"bytes"
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

type EsClient struct {
	es        *elasticsearch.Client
	bi        esutil.BulkIndexer
	processed uint64
}

func NewEsClient(dsn string) *EsClient {
	retryBackoff := backoff.NewExponentialBackOff()
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Username:      "elastic",
		Password:      "123qQ123",
		Addresses:     []string{dsn},
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 5,
	})
	if err != nil {
		panic(err)
	}
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         "prices",         // The default index name
		Client:        es,               // The Elasticsearch client
		NumWorkers:    3,                // The number of worker goroutines
		FlushBytes:    int(4096),        // The flush threshold in bytes
		FlushInterval: 10 * time.Second, // The periodic flush interval
	})
	if err != nil {
		panic(err)
	}
	return &EsClient{es: es, bi: bi}
}

func (e *EsClient) Ping() error {
	_, err := e.es.Ping()
	return err
}

func (e *EsClient) Insert(ctx context.Context, records []*Record) error {
	for _, record := range records {
		data, _ := sonic.Marshal(record)
		if err := e.bi.Add(ctx, esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: record.UUID,
			Body:       bytes.NewReader(data),
			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				atomic.AddUint64(&e.processed, 1)
			},
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("ERROR: %s", err)
				} else {
					log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
				}
			},
		}); err != nil {
			log.Println("error:", err)
		}
	}
	return nil
}
