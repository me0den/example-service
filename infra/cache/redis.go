package cache

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/me0den/example-service/infra/config"
	mRedis "github.com/me0den/example-service/x/redis"
)

var RedisFXModule = fx.Provide(
	NewRedis,
)

func NewRedis(cfg *config.Config) (*redis.Client, error) {
	client, err := mRedis.New(&cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("error when init redis client: %v", err)
	}

	return client.Client(), nil
}
