package responses

type SignMessage struct {
	PubKey    string `json:"pub_key"`
	Signature string `json:"signature"`
}
