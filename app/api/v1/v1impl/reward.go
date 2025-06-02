package v1impl

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	v1 "github.com/me0den/example-service/app/api/v1"
	"github.com/me0den/example-service/domain/entity"
	"github.com/me0den/example-service/domain/repo"
)

type RewardService struct {
	redisRepo repo.RedisRepo
}

func NewRewardService(
	redisRepo repo.RedisRepo,
) v1.RewardService {
	svc := &RewardService{
		redisRepo: redisRepo,
	}

	return svc
}

func (s *RewardService) CreateReward(c echo.Context) error {
	req := new(v1.CreateRewardRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	userElos, err := s.listUserElos(ctx, req.Teams)
	if err != nil {
		return err
	}

	winnerIndex := req.GetWinnerIndex()
	newUserElos := s.calculateElo(ctx, userElos, winnerIndex)
	res := &v1.CreateRewardResponse{}
	for idx, elo := range newUserElos {
		rankReward := &v1.Reward{
			NewElo: elo.Elo,
			OldElo: userElos[idx].Elo,
			UserID: elo.UserID,
		}

		res.Items = append(res.Items, rankReward)
	}

	if err := s.redisRepo.BatchUpdateElo(ctx, newUserElos); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &res)
}

func (s *RewardService) listUserElos(ctx context.Context, teams []*entity.Team) ([]*entity.UserElo, error) {
	var userElos []*entity.UserElo
	for _, team := range teams {
		userElo, err := s.redisRepo.GetUserElo(ctx, team.Owner)
		if err != nil {
			return nil, err
		}

		userElos = append(userElos, userElo)
	}

	return userElos, nil
}

func (s *RewardService) calculateElo(ctx context.Context, userElos []*entity.UserElo, winnerIdx int) []*entity.UserElo {
	newUserElos := []*entity.UserElo{
		userElos[0].Clone(),
		userElos[1].Clone(),
	}
	switch winnerIdx {
	case 0:
		newUserElos[0].Elo += 5
		newUserElos[1].Elo += 5
	case 1:
		newUserElos[0].Elo += 10
		newUserElos[1].Elo -= 10
	case 2:
		newUserElos[1].Elo += 10
		newUserElos[0].Elo -= 10
	}

	return newUserElos
}
