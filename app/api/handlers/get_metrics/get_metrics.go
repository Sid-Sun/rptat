package get_metrics

import (
	"encoding/json"
	"github.com/sid-sun/rptat/app/proxy/metrics"
	"github.com/sid-sun/rptat/app/proxy/service"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Handler handles all get data requests
func Handler(svc *service.Service, mtr *metrics.Metrics, lgr *zap.Logger, count int) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		mtr.SyncNow()

		met := (*svc).GetCurrentMetrics()

		var raw []byte
		var err error
		if count == len(met) {
			raw, err = json.Marshal(met)
		} else {
			raw, err = json.Marshal(getSubMetrics(met, count))
		}

		if err != nil {
			lgr.Sugar().Errorf("[Handlers] [GetHandler] [Marshal]")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = writer.Write(raw)
	}
}

func getSubMetrics(m service.Metrics, count int) *service.Metrics {
	subMetrics := make(service.Metrics)
	for i := 0; i < count && i < len(m); i++ {
		date := time.Now().AddDate(0, 0, -i).Format("01-02-2006")
		subMetrics[date] = m[date]
	}
	return &subMetrics
}
