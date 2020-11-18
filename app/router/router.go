package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sid-sun/rptat/app/api/handlers"
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
	rtr.Get("/rptat/api/getall", pxy.Authenticator.JustCheck(middlewares.WithContentJSON(handlers.
		GetHandler(pxy.Service, pxy.Metrics, lgr))))

	return rtr
}
