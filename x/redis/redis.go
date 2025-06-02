package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	TLSConfig *struct {
		CertFilePath       string `mapstructure:"cert_file_path"`
		InsecureSkipVerify bool   `mapstructure:"insecure_skip_verify"`
	} `mapstructure:"tls_config,omitempty"`
	Host         string        `mapstructure:"host"`
	Password     string        `mapstructure:"pass"`
	Port         int           `mapstructure:"port"`
	Database     int           `mapstructure:"database"`
	TTL          time.Duration `mapstructure:"ttl"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	WriteTimeOut time.Duration `mapstructure:"write_timeout"`
	ReadTimeOut  time.Duration `mapstructure:"read_timeout"`
	DialTimeOut  time.Duration `mapstructure:"dial_timeout"`
}

type Redis interface {
	Client() *redis.Client
}

type mRedis struct {
	client *redis.Client
}

func New(cfg *Config) (Redis, error) {
	if cfg.MinIdleConns == 0 {
		cfg.MinIdleConns = 10
	}
	redisOpt := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.Database,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		WriteTimeout: cfg.WriteTimeOut,
		ReadTimeout:  cfg.ReadTimeOut,
		DialTimeout:  cfg.DialTimeOut,
	}

	if cfg.TLSConfig != nil {
		tlsConfig := &tls.Config{
			//nolint:gosec
			InsecureSkipVerify: cfg.TLSConfig.InsecureSkipVerify,
			MinVersion:         tls.VersionTLS12,
		}

		if !cfg.TLSConfig.InsecureSkipVerify {
			caCert, err := os.ReadFile(cfg.TLSConfig.CertFilePath)
			if err != nil {
				return nil, err
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = caCertPool
			redisOpt.TLSConfig = tlsConfig
		}
	}

	// Connect to redis server
	client := redis.NewClient(redisOpt)
	log.Println("Pinging to Redis Server: ", cfg.Host, cfg.Port, cfg.Database)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Redis Server")
	return &mRedis{
		client: client,
	}, nil
}

func (r *mRedis) Client() *redis.Client {
	return r.client
}
