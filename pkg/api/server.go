package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sid-sun/rptat/cmd/config"
	"github.com/sid-sun/rptat/pkg/api/metrics"
	"github.com/sid-sun/rptat/pkg/api/proxy"
	"github.com/sid-sun/rptat/pkg/api/service"
	"github.com/sid-sun/rptat/pkg/api/store"
	"go.uber.org/zap"
)

// StartServer starts the api, inits all the requited submodules and routine for shutdown
func StartServer(cfg config.Config, logger *zap.Logger) {
	str := store.NewStore(cfg.StoreConfig, logger)
	svc := service.NewService(&str, logger)

	mtr, sync, err := metrics.NewMetrics(&svc)
	if err != nil {
		panic(err)
	}

	pxy, err := proxy.NewProxy("http://localhost:8081", logger, mtr)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", pxy.MetricsProxyHandler())

	srv := &http.Server{Addr: cfg.App.Address()}

	logger.Info(fmt.Sprintf("[StartServer] Listening on %s", cfg.App.Address()))
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error(fmt.Sprintf("[StartServer] [ListenAndServe]: %s", err.Error()))
			panic(err)
		}
	}()

	go mtr.Sync()
	gracefulShutdown(srv, logger, sync)
}

func gracefulShutdown(srv *http.Server, logger *zap.Logger, syncMetrics *chan bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logger.Info("Attempting GracefulShutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("[GracefulShutdown] [Shutdown]: %s", err.Error()))
			panic(err)
		}
	}()

	*syncMetrics <- false
	// Wait for ack for sync completion
	<-*syncMetrics
}
