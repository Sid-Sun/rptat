package main

import (
	"github.com/sid-sun/rptat/cmd/config"
	"github.com/sid-sun/rptat/pkg/api"
)

func main() {
	cfg := config.Load()
	initLogger(cfg.GetEnv())
	api.StartServer(cfg, logger)
}
