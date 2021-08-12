package rest

import (
	"encoding/json"
	"net/http"

	resttypes "github.com/sentinel-official/cli-client/types/rest"
)

func write(w http.ResponseWriter, status int, resp *resttypes.Response) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(resp)
}

func WriteErrorToResponse(w http.ResponseWriter, status int, result *resttypes.Error) {
	_ = write(
		w,
		status,
		resttypes.NewResponse(
			result,
			nil,
		),
	)
}

func WriteResultToResponse(w http.ResponseWriter, status int, result interface{}) {
	_ = write(
		w,
		status,
		resttypes.NewResponse(
			nil,
			result,
		),
	)
}
