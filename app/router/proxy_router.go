package router

import (
	auth "github.com/abbot/go-http-auth"
	"github.com/gorilla/mux"
	"github.com/sid-sun/rptat/app/api/handlers"
	"github.com/sid-sun/rptat/app/proxy"
	"go.uber.org/zap"
	"net/http"
)

func NewProxyRouter(proxies []proxy.Proxy, lgr *zap.Logger) *mux.Router {
	rtr := mux.NewRouter()

	for _, pxy := range proxies {
		authy := auth.NewDigestAuthenticator(pxy.AuthConfig.GetRealm(), auth.HtdigestFileProvider(pxy.AuthConfig.GetDigestFileName()))

		rtr.Handle("/", pxy.Handler.MetricsProxyHandler()).Host(pxy.Hostname)
		rtr.Handle("/getall", authy.JustCheck(withContentJSON(handlers.GetHandler(pxy.Service, pxy.Metrics, lgr)))).Methods(http.MethodGet).Host(pxy.Hostname)
	}

	return rtr
}

func withContentJSON(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	}
}
