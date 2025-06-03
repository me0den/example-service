package v1impl

import "go.uber.org/fx"

// FXModule represents a FX module for app api service.
var FXModule = fx.Provide(
	NewRewardService,
)
