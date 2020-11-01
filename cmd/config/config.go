package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config contains all the necessary configurations
type Config struct {
	App           appConfig
	environment   string
	StoreConfig   StoreConfig
	ProxyConfig   ProxyConfig
	MetricsConfig MetricsConfig
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

	fmt.Println(viper.GetString("APP_ENV"))
	return Config{
		environment: viper.GetString("APP_ENV"),
		App: appConfig{
			port: viper.GetString("APP_PORT"),
		},
		StoreConfig: StoreConfig{
			fileName:  viper.GetString("FILE_NAME"),
			filePerms: viper.GetInt("FILE_PERMS"),
		},
		ProxyConfig: ProxyConfig{
			listen: listenCfg{
				host: viper.GetString("PROXY_LISTEN_HOST"),
				port: viper.GetInt("PROXY_LISTEN_PORT"),
			},
			serve: serveCfg{
				protocol: viper.GetString("PROXY_SERVE_PROTOCOL"),
				host:     viper.GetString("PROXY_SERVE_HOST"),
				port:     viper.GetInt("PROXY_SERVE_PORT"),
			},
		},
		MetricsConfig: MetricsConfig{
			minForSync:           viper.GetUint("METRICS_MAX_PENDING"),
			periodicSyncInterval: viper.GetUint("METRICS_PERIODIC_SYNC_INTERVAL_IN_SECONDS"),
		},
	}
}
