package keys

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sentinel-official/cli-client/context"
)

func RegisterRoutes(r *mux.Router, ctx *context.Context) {
	r.Name("GetKey").
		Methods(http.MethodGet).Path("/keys/{name}").
		HandlerFunc(HandlerGetKey(ctx))
	r.Name("GetKeys").
		Methods(http.MethodGet).Path("/keys").
		HandlerFunc(HandlerGetKeys(ctx))
	r.Name("AddKey").
		Methods(http.MethodPost).Path("/keys").
		HandlerFunc(HandlerAddKey(ctx))
	r.Name("SignTx").
		Methods(http.MethodPost).Path("/keys/{name}/sign").
		HandlerFunc(HandlerSignTx(ctx))
	r.Name("DeleteKey").
		Methods(http.MethodDelete).Path("/keys/{name}").
		HandlerFunc(HandlerDeleteKey(ctx))
}
