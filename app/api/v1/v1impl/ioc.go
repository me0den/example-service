package v1impl

import "go.uber.org/fx"

var FXModule = fx.Provide(
	NewRewardService,
)
