package service

type ResponseStatus struct {
	From     string `json:"from"`
	ID       uint64 `json:"id"`
	To       string `json:"to"`
	Upload   int64  `json:"upload"`
	Download int64  `json:"download"`
}
