package proxy_router

import (
	"github.com/gorilla/mux"
	"github.com/sid-sun/rptat/app/proxy"
	"net/http"
)

func NewProxyRouter(proxies []*proxy.Proxy) *mux.Router {
	rtr := mux.NewRouter()

	for _, pxy := range proxies {
		rtr.HandleFunc("/", pxy.MetricsProxyHandler()).Host(pxy.Hostname)
		//rtr.HandleFunc("/", pxy.MetricsProxyHandler())
	}

	return rtr
}

func WithContentJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}
