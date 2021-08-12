package responses

type SignBytes struct {
	PubKey    string `json:"pub_key"`
	Signature string `json:"signature"`
}
