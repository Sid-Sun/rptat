package main

import (
	"github.com/sid-sun/rptat/cmd/config"
	"github.com/sid-sun/rptat/app"
)

func main() {
	cfg := config.Load()
	initLogger(cfg.GetEnv())
	app.StartServer(cfg, logger)
}
