package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sid-sun/rptat/app/api/handlers"
	"github.com/sid-sun/rptat/app/service"
	"go.uber.org/zap"
)

// NewRouter returns a new router instance
func NewRouter(svc *service.Service, lgr *zap.Logger) *mux.Router {
	rtr := mux.NewRouter()

	rtr.Handle("/getall", handlers.GetHandler(svc, lgr)).Methods(http.MethodPost)

	return rtr
}
