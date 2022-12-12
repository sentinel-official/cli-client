package responses

type GenerateSignature struct {
	PubKey    string `json:"pub_key"`
	Signature string `json:"signature"`
}
