package proxy

import (
	auth "github.com/abbot/go-http-auth"
	"github.com/sid-sun/rptat/app/proxy/metrics"
	"github.com/sid-sun/rptat/app/proxy/proxy_handler"
	"github.com/sid-sun/rptat/app/proxy/service"
	"github.com/sid-sun/rptat/app/proxy/store"
	"github.com/sid-sun/rptat/cmd/config"
	"go.uber.org/zap"
)

// NewProxy instantiates all the requirements for a proxy and returns them bundles in Proxy
func NewProxy(cfg *config.ProxyConfig, lgr *zap.Logger) Proxy {
	str := store.NewStore(cfg.StoreConfig, lgr)

	svc, err := service.NewService(&str, lgr)
	if err != nil {
		panic(err)
	}

	mtr, err := metrics.NewMetrics(&svc, cfg.MetricsConfig)
	if err != nil {
		panic(err)
	}

	p, err := proxy_handler.NewProxyHandler(*cfg, lgr, mtr)
	if err != nil {
		panic(err)
	}

	return Proxy{
		Handler:       p,
		Service:       &svc,
		Metrics:       mtr,
		Hostname:      cfg.GetHostname(),
		Authenticator: auth.NewDigestAuthenticator(cfg.AuthConfig.GetRealm(), auth.HtdigestFileProvider(cfg.AuthConfig.GetDigestFileName())),
	}
}
