package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sid-sun/rptat/app/metrics"
	"github.com/sid-sun/rptat/app/proxy"
	"github.com/sid-sun/rptat/app/service"
	"github.com/sid-sun/rptat/app/store"
	"github.com/sid-sun/rptat/cmd/config"
	"go.uber.org/zap"
)

// StartServer starts the proxy, inits all the requited submodules and routine for shutdown
func StartServer(cfg config.Config, logger *zap.Logger) {
	str := store.NewStore(cfg.StoreConfig, logger)

	svc, err := service.NewService(&str, logger)
	if err != nil {
		panic(err)
	}

	mtr, sync, err := metrics.NewMetrics(&svc, cfg.MetricsConfig)
	if err != nil {
		panic(err)
	}

	pxy, err := proxy.NewProxy(cfg.ProxyConfig, logger, mtr)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", pxy.MetricsProxyHandler())

	srv := &http.Server{Addr: cfg.ProxyConfig.GetListenAddress()}

	logger.Info(fmt.Sprintf("[StartServer] Listening on %s", cfg.App.Address()))
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
