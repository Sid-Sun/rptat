package main

import (
	"github.com/sid-sun/rptat/app"
	"github.com/sid-sun/rptat/cmd/config"
)

func main() {
	cfg := config.Load()
	initLogger(cfg.GetEnv())
	app.StartServer(cfg, logger)
}
