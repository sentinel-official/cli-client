package rest

type Response struct {
	Success bool        `json:"success"`
	Error   *Error      `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

func NewResponse(e *Error, result interface{}) *Response {
	return &Response{
		Success: e == nil,
		Error:   e,
		Result:  result,
	}
}
