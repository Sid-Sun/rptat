package app

import (
	"context"
	"fmt"
	"github.com/sid-sun/rptat/app/api/router"
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
	str := store.NewStore(cfg.ProxyConfig.Store, logger)

	svc, err := service.NewService(&str, logger)
	if err != nil {
		panic(err)
	}

	mtr, err := metrics.NewMetrics(&svc, cfg.ProxyConfig.Metrics)
	if err != nil {
		panic(err)
	}

	pxy, err := proxy.NewProxy(&cfg.ProxyConfig, logger, mtr)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", pxy.MetricsProxyHandler())

	proxyServer := &http.Server{Addr: cfg.ProxyConfig.GetListenAddress()}

	logger.Info(fmt.Sprintf("[StartServer] [Proxy] Listening on %s", cfg.ProxyConfig.GetListenAddress()))
	go func() {
		if err := proxyServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("[StartServer] [Proxy] [ListenAndServe]: %s", err.Error()))
			panic(err)
		}
	}()

	rtr := router.NewRouter(&svc, mtr, logger)
	apiServer := &http.Server{Addr: cfg.API.Address(), Handler: rtr}

	logger.Info(fmt.Sprintf("[StartServer] [API] Listening on %s", cfg.API.Address()))
	go func() {
		if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("[StartServer] [API] [ListenAndServe]: %s", err.Error()))
			panic(err)
		}
	}()

	gracefulShutdown(apiServer, proxyServer, logger, mtr)
}

func gracefulShutdown(apiServer, proxyServer *http.Server, logger *zap.Logger, mtr *metrics.Metrics) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logger.Info("Attempting GracefulShutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := proxyServer.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("[GracefulShutdown] [Proxy] [Shutdown]: %s", err.Error()))
			panic(err)
		}
	}()

	go func() {
		if err := apiServer.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("[GracefulShutdown] [API] [Shutdown]: %s", err.Error()))
			panic(err)
		}
	}()

	// Perform a blocking sync
	mtr.SyncAndShutdown()
}
