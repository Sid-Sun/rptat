package config

import "fmt"

type apiConfig struct {
	host string
	port string
}

// Address returns the requisite address router should listen at
func (ac apiConfig) Address() string {
	return fmt.Sprintf("%s:%s", ac.host, ac.port)
}
