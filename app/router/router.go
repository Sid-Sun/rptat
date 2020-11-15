package router

import (
	"github.com/gorilla/mux"
	"github.com/sid-sun/rptat/app/api/handlers"
	"github.com/sid-sun/rptat/app/api/middlewares"
	"github.com/sid-sun/rptat/app/proxy"
	"go.uber.org/zap"
	"net/http"
)

func NewProxyRouter(proxies []proxy.Proxy, lgr *zap.Logger) *mux.Router {
	rtr := mux.NewRouter()

	for _, pxy := range proxies {
		rtr.Handle("/", pxy.Handler.MetricsProxyHandler()).Host(pxy.Hostname)
		rtr.Handle("/rptat/api", pxy.Authenticator.JustCheck(middlewares.WithContentJSON(handlers.
			GetHandler(pxy.Service, pxy.Metrics, lgr)))).Methods(http.MethodGet).Host(pxy.Hostname)
	}

	return rtr
}
