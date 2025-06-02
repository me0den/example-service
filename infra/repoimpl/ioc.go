package repoimpl

import (
	"go.uber.org/fx"
)

var FXModule = fx.Provide(
	NewRedisDBRepo,
)
