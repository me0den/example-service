package config

import (
	"errors"

	"github.com/spf13/viper"
	"go.uber.org/fx"

	"github.com/me0den/example-service/x/redis"
)

// Config is a group of options for the service.
type Config struct {
	HTTPServer struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"http_server"`
	Redis redis.Config `mapstructure:"redis"`
}

// Load loads Config from Viper and returns them.
func Load(v *viper.Viper) (*Config, error) {
	cfg := &Config{}
	if err := v.Unmarshal(&cfg); err != nil {
		return &Config{}, errors.New(err.Error())
	}

	return cfg, nil
}

// FXModule represents a FX module for config.
var FXModule = fx.Options(
	fx.Provide(
		Load,
	),
)
