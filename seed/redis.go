package dbench

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
}

func NewRedis() *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &Redis{Client: rdb}
}

func (r *Redis) Insert(ctx context.Context, records []*Record) error {
	pipe := r.Client.Pipeline()
	for _, record := range records {
		data, _ := sonic.MarshalString(record)
		pipe.Set(ctx, record.UUID, data, 0)
	}
	_, err := pipe.Exec(ctx)
	return err
}
