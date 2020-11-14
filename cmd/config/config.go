package config

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"github.com/sid-sun/rptat/cmd/config/internal"
	"io/ioutil"
	"os"
)

// Config contains all the necessary configurations
type Config struct {
	API           apiConfig
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
		API: apiConfig{
			host: co.API.Host,
			port: co.API.Port,
		},
		StoreConfig: StoreConfig{
			fileName:  co.StoreConfig.FileName,
			filePerms: co.StoreConfig.FilePerms,
		},
		ProxyConfig: ProxyConfig{
			protocol: co.ProxyConfig.Protocol,
			port:     co.ProxyConfig.Port,
			host:     co.ProxyConfig.Host,
			hostname: co.ProxyConfig.Hostname,
		},
		MetricsConfig: MetricsConfig{
			minForSync:           co.MetricsConfig.MinForSync,
			periodicSyncInterval: co.MetricsConfig.PeriodicSyncInterval,
		},
	}

	d, err = toml.Marshal(co)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("config.toml", d, 420)
	if err != nil {
		panic(err)
	}

	fmt.Println(c)
	return c
}
