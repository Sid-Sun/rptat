package config

import "fmt"

type (
	// ProxyConfig defines config for proxies
	ProxyConfig struct {
		protocol string
		port     int
		host     string
		hostname string
		Store    *StoreConfig
		Metrics  *MetricsConfig
	}
)

// GetServeURL defines the URL for resource to be proxied
func (p ProxyConfig) GetServeURL() string {
	return fmt.Sprintf("%s://%s:%d", p.protocol, p.host, p.port)
}

// GetHostname defines the URL for resource to be proxied
func (p ProxyConfig) GetHostname() string {
	return p.hostname
}
