package responses

type GetStatus struct {
	ID       uint64 `json:"id"`
	IFace    string `json:"iface"`
	Upload   int64  `json:"upload"`
	Download int64  `json:"download"`
}
