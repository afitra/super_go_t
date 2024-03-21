package connection

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"superindo/v1/model"
)

type Redis_connection struct {
	redis_client *redis.Client
}

func NewRedis_Connection(redis *redis.Client) Redis_cache {
	return &Redis_connection{redis}
}

func (IN *Redis_connection) Redis_set_key(ctx context.Context, payload model.Chache_model) error {
	var err error

	if err = IN.redis_client.Set(ctx, payload.Key, string(payload.Data), payload.Expired).Err(); err != nil {
		return err
	}

	return nil
}

func (IN *Redis_connection) Redis_get_key(ctx context.Context, key string) (interface{}, error) {
	var err error
	var cache_data string
	var data interface{}

	if cache_data, err = IN.redis_client.Get(ctx, key).Result(); err != nil {
		return nil, nil
	}

	if err := json.Unmarshal([]byte(cache_data), &data); err != nil {
		return nil, err
	}

	return data, nil
}
