package config

import (
	"github.com/pelletier/go-toml"
	"github.com/sid-sun/rptat/cmd/config/internal"
	"io/ioutil"
	"os"
)

// Config contains all the necessary configurations
type Config struct {
	API         appConfig
	environment string
	ProxyConfig []ProxyConfig
}

// GetEnv returns the current environment
func (c Config) GetEnv() string {
	return c.environment
}

// Load reads all config from env to config
func Load() Config {
	f, err := os.Open("config.toml")
	if err != nil {
		panic(err)
	}
	d, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	co := internal.Config{}
	err = toml.Unmarshal(d, &co)
	if err != nil {
		panic(err)
	}

	c := Config{
		environment: co.App.Env,
		API: appConfig{
			host: co.App.Host,
			port: co.App.Port,
		},
		ProxyConfig: *new([]ProxyConfig),
	}

	for _, pxy := range co.ProxyConfig {
		c.ProxyConfig = append(c.ProxyConfig, ProxyConfig{
			protocol: pxy.Protocol,
			port:     pxy.Port,
			host:     pxy.Host,
			hostname: pxy.Hostname,
			StoreConfig: &StoreConfig{
				fileName:  pxy.StoreConfig.FileName,
				filePerms: pxy.StoreConfig.FilePerms,
			},
			MetricsConfig: &MetricsConfig{
				minForSync:           pxy.MetricsConfig.MinForSync,
				periodicSyncInterval: pxy.MetricsConfig.PeriodicSyncInterval,
			},
			AuthConfig: &Auth{
				htDigestFile: pxy.AuthConfig.HTDigestFile,
				realm:        pxy.AuthConfig.Realm,
			},
		})
	}

	d, err = toml.Marshal(co)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("config.toml", d, 420)
	if err != nil {
		panic(err)
	}

	return c
}
