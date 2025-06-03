package v1

import (
	"github.com/labstack/echo/v4"

	"github.com/me0den/example-service/domain/entity"
)

type RewardService interface {
	CreateReward(c echo.Context) error
}

type Reward struct {
	UserID    string `json:"userID"`
	OldElo    int    `json:"oldElo"`
	NewElo    int    `json:"newElo"`
	UpdatedAt int64  `json:"updatedAt"`
}

type Rewards struct {
	Items []*Reward `json:"rewards"`
}

type CreateRewardRequest struct {
	Winner string         `json:"winner" validate:"required"`
	Teams  []*entity.Team `json:"teams" validate:"required,eq=2"`
}

// GetWinnerIndex retrieve index of winner from request
//
// 0 means draw
//
// 1 means first team wins
//
// 2 means second team wins
func (c *CreateRewardRequest) GetWinnerIndex() int {
	winnerIndex := 0
	if c.Winner == c.Teams[0].Owner {
		winnerIndex = 1
	} else if c.Winner == c.Teams[1].Owner {
		winnerIndex = 2
	}

	return winnerIndex
}

type CreateRewardResponse = Rewards
