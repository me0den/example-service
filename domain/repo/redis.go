package repo

import (
	"context"

	"github.com/me0den/example-service/domain/entity"
)

// RedisRepo provides methods for interacting with redis data.
type RedisRepo interface {
	GetUserElo(ctx context.Context, userID string) (*entity.UserElo, error)
	BatchUpdateElo(ctx context.Context, elos []*entity.UserElo) error
}
