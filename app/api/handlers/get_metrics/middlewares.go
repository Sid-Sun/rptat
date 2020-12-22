package get_metrics

import (
	"github.com/sid-sun/rptat/app/proxy/metrics"
	"github.com/sid-sun/rptat/app/proxy/service"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func HandlerWithCount(svc *service.Service, mtr *metrics.Metrics, lgr *zap.Logger) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		c := request.URL.Query().Get("count")
		if c == "" {
			panic("Count value empty")
		}

		count, err := strconv.Atoi(c)
		if err != nil {
			panic(err)
		}
		if count <= 0 {
			panic("Count cannot be zero or non-negative")
		}

		Handler(svc, mtr, lgr, count).ServeHTTP(writer, request)
	})
}
