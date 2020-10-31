package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sid-sun/rptat/app/service"
	"go.uber.org/zap"
)

// GetHandler handles all get data requests
func GetHandler(svc *service.Service, lgr *zap.Logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		raw, err := json.Marshal((*svc).GetCurrentMetrics())
		if err != nil {
			lgr.Sugar().Errorf("[Handlers] [GetHandler] [Marshal]")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Add("Content-Type", "application/json")
		_, err = writer.Write(raw)

	}
}
