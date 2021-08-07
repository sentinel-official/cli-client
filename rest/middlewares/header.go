package middlewares

import (
	"net/http"

	"github.com/go-kit/kit/transport/http/jsonrpc"
)

func AddHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", jsonrpc.ContentType)
			next.ServeHTTP(w, r)
		},
	)
}
