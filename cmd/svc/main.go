package main

import (
	"go.uber.org/fx"

	"github.com/me0den/example-service/app/api/v1/transport/routes"
	"github.com/me0den/example-service/app/api/v1/v1impl"
	"github.com/me0den/example-service/infra/cache"
	"github.com/me0den/example-service/infra/config"
	"github.com/me0den/example-service/infra/repoimpl"
	"github.com/me0den/example-service/x/viper"
)

func main() {
	app := fx.New(
		viper.FXModule,
		config.FXModule,
		routes.ServerFXModule,
		cache.RedisFXModule,
		repoimpl.FXModule,
		v1impl.FXModule,
	)
	app.Run()
}
