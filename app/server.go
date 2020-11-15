package app

import (
	"context"
	"fmt"
	"github.com/sid-sun/rptat/app/proxy"
	"github.com/sid-sun/rptat/app/router"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sid-sun/rptat/cmd/config"
	"go.uber.org/zap"
)

// StartServer starts the proxy, inits all the requited submodules and routine for shutdown
func StartServer(cfg config.Config, logger *zap.Logger) {
	proxies := *new([]proxy.Proxy)

	for _, pxy := range cfg.ProxyConfig {
		proxies = append(proxies, proxy.NewProxy(&pxy, logger))
		logger.Sugar().Infof("Subscribed [%s] as [%s]", pxy.GetServeURL(), pxy.GetHostname())
	}

	proxyRouter := router.NewProxyRouter(proxies, logger)
	proxyServer := &http.Server{Addr: cfg.API.Address(), Handler: proxyRouter}

	logger.Info(fmt.Sprintf("[StartServer] [Server] Listening on %s", cfg.API.Address()))
	go func() {
		if err := proxyServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("[StartServer] [Server] [ListenAndServe]: %s", err.Error()))
			panic(err)
		}
	}()

	gracefulShutdown(proxyServer, logger, proxies)
}

func gracefulShutdown(httpServer *http.Server, logger *zap.Logger, pxy []proxy.Proxy) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logger.Info("Attempting GracefulShutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := httpServer.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("[GracefulShutdown] [Server] [Shutdown]: %s", err.Error()))
			panic(err)
		}
	}()

	for _, pxy := range pxy {
		pxy.Metrics.SyncAndShutdown()
	}
}
