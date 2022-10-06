package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	MongoDBUri string `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri   string `mapstructure:"REDIS_URI"`
	ServerPort string `mapstructure:"SERVER_PORT"`
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

func(cfg *Config) validateConfigs() error {
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
