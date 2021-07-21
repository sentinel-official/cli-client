package keyring

import (
	"encoding/hex"

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
		PubKey:  hex.EncodeToString(v.GetPubKey().Bytes()),
		Address: hex.EncodeToString(v.GetAddress().Bytes()),
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
