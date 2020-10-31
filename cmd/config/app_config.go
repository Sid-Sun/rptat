package config

import "fmt"

type appConfig struct {
	port string
}

// Address returns the requisite address router should listen at
func (ac appConfig) Address() string {
	return fmt.Sprintf("localhost:%s", ac.port)
}
