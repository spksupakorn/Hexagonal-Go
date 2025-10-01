package config

import (
	"sync"

	"github.com/spf13/viper"
)

var (
	once sync.Once
	vipe = viper.New()
)

// LoadConfig loads configuration from .env file and environment variables (singleton)
func LoadConfig() {
	once.Do(func() {
		vipe.SetConfigFile(".env")
		vipe.AutomaticEnv()
		vipe.ReadInConfig()
	})
}

func GetConfigString(key string) string {
	return vipe.GetString(key)
}

func GetConfigInt(key string) int {
	return vipe.GetInt(key)
}

func GetConfigInt64(key string) int64 {
	return vipe.GetInt64(key)
}
