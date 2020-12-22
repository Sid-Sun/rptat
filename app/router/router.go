package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sid-sun/rptat/app/api/handlers/get_metrics"
	"github.com/sid-sun/rptat/app/api/middlewares"
	"github.com/sid-sun/rptat/app/proxy"
	"go.uber.org/zap"
)

func NewRouter() *chi.Mux {
	rtr := chi.NewRouter()
	rtr.Use(middleware.Recoverer)
	return rtr
}

func NewProxyRouter(pxy proxy.Proxy, lgr *zap.Logger) chi.Router {
	rtr := chi.NewRouter()

	rtr.Handle("/*", pxy.Handler.MetricsProxyHandler())
	rtr.Get("/rptat/api/get/all", middlewares.WithContentJSON(get_metrics.
		Handler(pxy.Service, pxy.Metrics, lgr, len((*pxy.Service).GetCurrentMetrics()))))
	rtr.Get("/rptat/api/get", middlewares.WithContentJSON(get_metrics.
		HandlerWithCount(pxy.Service, pxy.Metrics, lgr)))

	return rtr
}
