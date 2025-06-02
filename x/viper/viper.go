package viper

import (
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	envPrefix        = "SVC"
	envConfigEnabled = "ENV_CONFIG_ENABLED"
	envConfigPath    = "CONFIG_PATH"
)

var FXModule = fx.Options(
	fx.Provide(
		NewViper,
	),
)

func NewViper() (*viper.Viper, error) {
	configPath := os.Getenv(envConfigPath)
	if configPath == "" {
		configPath = "./infra/config"
	}
	return NewViperFrom(configPath)
}

func NewViperFrom(path string) (*viper.Viper, error) {
	v := viper.New()

	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	isEnvConfig := os.Getenv(envConfigEnabled)
	if isEnvConfig == "" {
		isEnvConfig = "false"
	}

	envOnly, err := strconv.ParseBool(isEnvConfig)
	if err != nil {
		return nil, err
	}

	if !envOnly {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(path)

		return v, v.ReadInConfig()
	}

	return v, nil
}
