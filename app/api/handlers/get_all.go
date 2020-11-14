package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sid-sun/rptat/app/proxy/metrics"
	"github.com/sid-sun/rptat/app/proxy/service"
	"go.uber.org/zap"
)

// GetHandler handles all get data requests
func GetHandler(svc *service.Service, mtr *metrics.Metrics, lgr *zap.Logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		mtr.SyncNow()

		raw, err := json.Marshal((*svc).GetCurrentMetrics())
		if err != nil {
			lgr.Sugar().Errorf("[Handlers] [GetHandler] [Marshal]")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = writer.Write(raw)
	}
}
