package types

import (
	"net/http"
)

type RestError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewRestError(code int, message string) *RestError {
	return &RestError{
		Code:    code,
		Message: message,
	}
}

type RestResponse struct {
	Success bool        `json:"success"`
	Error   *RestError  `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

func NewRestResponse(err *RestError, res interface{}) *RestResponse {
	return &RestResponse{
		Success: err == nil,
		Error:   err,
		Result:  res,
	}
}

type RestResponseWriter struct {
	http.ResponseWriter
	Status int
	Length int
}

func NewRestResponseWriter(w http.ResponseWriter) *RestResponseWriter {
	return &RestResponseWriter{
		ResponseWriter: w,
		Status:         0,
		Length:         0,
	}
}

func (r *RestResponseWriter) Header() http.Header {
	return r.ResponseWriter.Header()
}

func (r *RestResponseWriter) Write(p []byte) (n int, err error) {
	n, err = r.ResponseWriter.Write(p)
	r.Length += n

	return n, err
}

func (r *RestResponseWriter) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.Status = status
}
