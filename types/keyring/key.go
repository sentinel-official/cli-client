package keyring

import (
	"encoding/base64"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

type Key struct {
	Name    string `json:"name"`
	PubKey  string `json:"pub_key"`
	Address string `json:"address"`
}

func NewKeyFromRaw(v keyring.Info) Key {
	return Key{
		Name:    v.GetName(),
		PubKey:  base64.StdEncoding.EncodeToString(v.GetPubKey().Bytes()),
		Address: base64.StdEncoding.EncodeToString(v.GetAddress().Bytes()),
	}
}

type (
	Keys []Key
)

func NewKeysFromRaw(v []keyring.Info) Keys {
	items := make(Keys, 0, len(v))
	for i := 0; i < len(v); i++ {
		items = append(items, NewKeyFromRaw(v[i]))
	}

	return items
}
