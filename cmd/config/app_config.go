package config

import "fmt"

type appConfig struct {
	host string
	port uint
}

// Address returns the requisite address router should listen at
func (ac appConfig) Address() string {
	return fmt.Sprintf("%s:%d", ac.host, ac.port)
}
