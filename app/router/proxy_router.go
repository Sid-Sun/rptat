package router

import (
	"github.com/gorilla/mux"
	"github.com/sid-sun/rptat/app/api/handlers"
	"github.com/sid-sun/rptat/app/proxy"
	"go.uber.org/zap"
	"net/http"
)

func NewProxyRouter(proxies []proxy.Proxy, lgr *zap.Logger) *mux.Router {
	rtr := mux.NewRouter()

	for _, pxy := range proxies {
		rtr.HandleFunc("/", pxy.Handler.MetricsProxyHandler()).Host(pxy.Hostname)
		rtr.HandleFunc("/getall", handlers.GetHandler(pxy.Service, pxy.Metrics, lgr)).Methods(http.MethodGet).Host(pxy.Hostname)
		//rtr.HandleFunc("/", pxy.Handler.MetricsProxyHandler())
	}

	return rtr
}

func WithContentJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}
