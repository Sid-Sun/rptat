package config

import "fmt"

type (
	ProxyConfig struct {
		listen listenCfg
		serve  serveCfg
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

func (p ProxyConfig) GetListenAddress() string {
	return fmt.Sprintf("%s:%d", p.listen.host, p.listen.port)
}

func (p ProxyConfig) GetServeURL() string {
	return fmt.Sprintf("%s://%s:%d", p.serve.protocol, p.serve.host, p.serve.port)
}
