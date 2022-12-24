package modules

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sentinel-official/cli-client/context"
	"github.com/sentinel-official/cli-client/rest/handlers"
	"github.com/sentinel-official/cli-client/rest/routes"
)

func RegisterKeyring(r *mux.Router, ctx *context.ServerContext) {
	r.Name(routes.GetKey).
		Methods(http.MethodPost).Path(routes.GetKey).
		HandlerFunc(handlers.GetKey(ctx))
	r.Name(routes.GetKeys).
		Methods(http.MethodPost).Path(routes.GetKeys).
		HandlerFunc(handlers.GetKeys(ctx))
	r.Name(routes.AddKey).
		Methods(http.MethodPost).Path(routes.AddKey).
		HandlerFunc(handlers.AddKey(ctx))
	r.Name(routes.SignMessage).
		Methods(http.MethodPost).Path(routes.SignMessage).
		HandlerFunc(handlers.SignMessage(ctx))
	r.Name(routes.Delete).
		Methods(http.MethodPost).Path(routes.Delete).
		HandlerFunc(handlers.DeleteKey(ctx))
}
