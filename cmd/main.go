package main

import (
	"github.com/sid-sun/rptat/app"
	"github.com/sid-sun/rptat/cmd/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	initLogger(cfg.GetEnv())
	app.StartServer(cfg, logger)
}
