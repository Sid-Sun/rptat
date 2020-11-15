package config

import "fmt"

type (
	// ProxyConfig defines config for proxies
	ProxyConfig struct {
		protocol      string
		port          int
		host          string
		hostname      string
		StoreConfig   *StoreConfig
		MetricsConfig *MetricsConfig
		AuthConfig    *Auth
	}
	// Auth defines the authentication config
	Auth struct {
		htDigestFile string
		realm        string
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

// GetDigestFileName returns the HTTP Digest file name
func (a *Auth) GetDigestFileName() string {
	return a.htDigestFile
}

// GetRealm returns thr realm for Digest auth
func (a *Auth) GetRealm() string {
	return a.realm
}
