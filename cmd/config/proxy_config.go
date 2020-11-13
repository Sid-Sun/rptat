package config

import "fmt"

type (
	// ProxyConfig defines config for proxies
	ProxyConfig struct {
		name string
		listen listenCfg
		serve  serveCfg
		Store StoreConfig
		Metrics MetricsConfig
	}
	listenCfg struct {
		port int
		host string
	}
	serveCfg struct {
		protocol string
		port     int
		host     string
	}
)

func (p ProxyConfig) GetName() string {
	return p.name
}

// GetListenAddress returns the address proxy should listen at
func (p ProxyConfig) GetListenAddress() string {
	return fmt.Sprintf("%s:%d", p.listen.host, p.listen.port)
}

// GetServeURL defines the URL for resource to be proxied
func (p ProxyConfig) GetServeURL() string {
	return fmt.Sprintf("%s://%s:%d", p.serve.protocol, p.serve.host, p.serve.port)
}
