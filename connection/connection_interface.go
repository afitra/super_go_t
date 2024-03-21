package connection

import (
	"context"
	"superindo/v1/model"
)

type Redis_cache interface {
	Redis_set_key(ctx context.Context, payload model.Chache_model) error
	Redis_get_key(ctx context.Context, key string) (interface{}, error)
}
