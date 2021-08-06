package modules

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sentinel-official/cli-client/context"
	"github.com/sentinel-official/cli-client/rest/handlers"
	"github.com/sentinel-official/cli-client/rest/routes"
)

func RegisterService(r *mux.Router, ctx *context.ServerContext) *mux.Router {
	r.Name(routes.Connect).
		Methods(http.MethodPost).Path(routes.Connect).
		Handler(handlers.Connect(ctx))
	r.Name(routes.Disconnect).
		Methods(http.MethodPost).Path(routes.Disconnect).
		Handler(handlers.Disconnect(ctx))
	r.Name(routes.GetStatus).
		Methods(http.MethodPost).Path(routes.GetStatus).
		Handler(handlers.GetStatus(ctx))

	return r
}
