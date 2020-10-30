package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config contains all the necessary configurations
type Config struct {
	App         appConfig
	environment string
	StoreConfig StoreConfig
}

// GetEnv returns the current environment
func (c Config) GetEnv() string {
	return c.environment
}

// Load reads all config from env to config
func Load() Config {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	fmt.Println(os.Getenv("APP_ENV"))
	return Config{
		environment: os.Getenv("APP_ENV"),
		App: appConfig{
			port: os.Getenv("APP_PORT"),
		},
		StoreConfig: StoreConfig{
			fileName:  os.Getenv("FILE_NAME"),
			filePerms: viper.GetInt("FILE_PERMS"),
		},
	}
}
