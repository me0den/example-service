package repoimpl

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/me0den/example-service/domain/entity"
	"github.com/me0den/example-service/domain/repo"
)

const (
	userEloKey = "user-elo"
)

type RedisRepo struct {
	client *redis.Client
}

// NewRedisDBRepo creates and returns a new instance of repo.RedisRepo.
func NewRedisDBRepo(
	client *redis.Client,
) repo.RedisRepo {
	return &RedisRepo{
		client: client,
	}
}

func (r *RedisRepo) GetUserElo(ctx context.Context, userID string) (*entity.UserElo, error) {
	data, err := r.client.HGet(ctx, userEloKey, userID).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	// Default elo: 1000
	userElo := &entity.UserElo{UserID: userID, Elo: 1000}
	if !errors.Is(err, redis.Nil) {
		if err := json.Unmarshal([]byte(data), userElo); err != nil {
			return nil, err
		}
	}

	return userElo, nil
}

func (r *RedisRepo) BatchUpdateElo(ctx context.Context, elos []*entity.UserElo) error {
	pipe := r.client.Pipeline()
	for _, elo := range elos {
		eloData, err := json.Marshal(elo)
		if err != nil {
			return err
		}

		pipe.HSet(ctx, userEloKey, elo.UserID, eloData)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	return nil
}
