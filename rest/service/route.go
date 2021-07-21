package service

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sentinel-official/cli-client/context"
)

func RegisterRoutes(r *mux.Router, ctx *context.Context) *mux.Router {
	r.Name("ServiceConnect").
		Methods(http.MethodPost).Path("/connect").
		Handler(HandlerConnect(ctx))
	r.Name("ServiceDisconnect").
		Methods(http.MethodPost).Path("/disconnect").
		Handler(HandlerDisconnect(ctx))
	r.Name("ServiceStatus").
		Methods(http.MethodGet).Path("/status").
		Handler(HandlerStatus(ctx))

	return r
}
