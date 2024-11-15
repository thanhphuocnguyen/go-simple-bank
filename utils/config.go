package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost                 string        `mapstructure:"DB_HOST"`
	DBPort                 int           `mapstructure:"DB_PORT"`
	DBUser                 string        `mapstructure:"DB_USER"`
	DBPassword             string        `mapstructure:"DB_PASSWORD"`
	DBName                 string        `mapstructure:"DB_NAME"`
	ServerAddress          string        `mapstructure:"SERVER_ADDRESS"`
	AccessTokenDuration    time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	SymmetricEncryptionKey string        `mapstructure:"SYMMETRIC_ENCRYPTION_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
