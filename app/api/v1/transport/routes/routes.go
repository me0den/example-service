package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"

	v1 "github.com/me0den/example-service/app/api/v1"
)

func RegisterRoutes(e *echo.Echo, rewardService v1.RewardService) {
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	groupV1 := e.Group("/v1")
	groupV1.POST("/battle/:battle_id/reward", rewardService.CreateReward)
}
