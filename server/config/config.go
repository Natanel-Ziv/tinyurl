package config

import "github.com/spf13/viper"

type Config struct {
	MongoDBUri string `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri   string `mapstructure:"REDIS_URI"`
	ServerPort string `mapstructure:"SERVER_PORT"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	return cfg, err
}
