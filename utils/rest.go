package utils

import (
	"encoding/json"
	"net/http"

	clitypes "github.com/sentinel-official/cli-client/types"
)

func write(w http.ResponseWriter, code int, body *clitypes.RestResponseBody) error {
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(body)
}

func WriteErrorToResponseBody(w http.ResponseWriter, code int, err *clitypes.RestError) {
	_ = write(
		w,
		code,
		clitypes.NewRestResponseBody(err, nil),
	)
}

func WriteResultToResponseBody(w http.ResponseWriter, code int, res interface{}) {
	_ = write(
		w,
		code,
		clitypes.NewRestResponseBody(nil, res),
	)
}
