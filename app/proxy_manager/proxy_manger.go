package proxy_manager

import (
	"context"
	"fmt"
	"github.com/sid-sun/rptat/app/metrics"
	"github.com/sid-sun/rptat/app/proxy"
	"github.com/sid-sun/rptat/app/service"
	"github.com/sid-sun/rptat/app/store"
	"github.com/sid-sun/rptat/cmd/config"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type prox struct {
	cfg *config.ProxyConfig
	met *metrics.Metrics
	svc *http.Server
}

// ProxyManager contains everything necessary for managing proxies
type ProxyManager struct {
	proxies map[string]prox
	lgr *zap.Logger
}

func NewProxyManager(config []*config.ProxyConfig, logger *zap.Logger) *ProxyManager {
	pm := new(ProxyManager)
	pm.proxies = *(new(map[string]prox))
	pm.lgr = logger
	for _, pxyCfg := range config {
		str := store.NewStore(pxyCfg.Store, logger)

		svc, err := service.NewService(&str, logger)
		if err != nil {
			panic(err)
		}

		mtr, err := metrics.NewMetrics(&svc, pxyCfg.Metrics)
		if err != nil {
			panic(err)
		}

		pxy, err := proxy.NewProxy(pxyCfg, logger, mtr)
		if err != nil {
			panic(err)
		}

		http.HandleFunc("/", pxy.MetricsProxyHandler())

		proxyServer := &http.Server{Addr: pxyCfg.GetListenAddress()}

		logger.Info(fmt.Sprintf("[StartServer] [Proxy] Listening on %s", pxyCfg.GetListenAddress()))
		go func() {
			if err := proxyServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error(fmt.Sprintf("[StartServer] [Proxy] [ListenAndServe]: %s", err.Error()))
				panic(err)
			}
		}()

		pm.proxies[pxyCfg.GetName()] = prox{
			met: mtr,
			svc: proxyServer,
			cfg: pxyCfg,
		}
	}
	return pm
}

func (p *ProxyManager) StopAll() {
	for name, pxy := range p.proxies {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		pxy := pxy
		go func() {
			if err := pxy.svc.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
				p.lgr.Error(fmt.Sprintf("[ProxyManager] [StopAll] %s [Shutdown]: %s", name, err.Error()))
				panic(err)
			}
		}()

		pxy.met.SyncAndShutdown()
		cancel()
	}
}
