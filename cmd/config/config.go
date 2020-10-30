package config

import (
	"fmt"
	"os"
)

// Config contains all the necessary configurations
type Config struct {
	App         appConfig
	environment string
}

// GetEnv returns the current environment
func (c Config) GetEnv() string {
	return c.environment
}

// Load reads all config from env to config
func Load() Config {
	fmt.Println(os.Getenv("APP_ENV"))
	return Config{
		environment: os.Getenv("APP_ENV"),
		App: appConfig{
			port: os.Getenv("APP_PORT"),
		},
	}
}
