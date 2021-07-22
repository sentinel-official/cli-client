package keys

type ResponseSignTx struct {
	PubKey    string `json:"pub_key"`
	Signature []byte `json:"signature"`
}
