package config

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	MongoDBUri             string        `mapstructure:"MONGODB_LOCAL_URI"`
	MongoDBName            string        `mapstructure:"MONGODB_DB_NAME"`
	RedisUri               string        `mapstructure:"REDIS_URI"`
	ServerPort             string        `mapstructure:"SERVER_PORT"`
	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("server.env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, cfg.validateConfigs()
}

func (cfg *Config) validateConfigs() error {
	if cfg.ServerPort == "" {
		return errors.New("must provide server port")
	}
	if cfg.MongoDBUri == "" {
		return errors.New("must provide Mongo DB URI")
	}
	if cfg.RedisUri == "" {
		return errors.New("must provide Redis URI")
	}
	return nil
}
