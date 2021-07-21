package keys

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmos/go-bip39"
	"github.com/pkg/errors"
)

type RequestAddKey struct {
	Name          string `json:"name"`
	Mnemonic      string `json:"mnemonic"`
	Type          uint32 `json:"type"`
	Account       uint32 `json:"account"`
	Index         uint32 `json:"index"`
	BIP39Password string `json:"bip39_password"`
}

func NewRequestAddKey(r *http.Request) (*RequestAddKey, error) {
	var v RequestAddKey
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (r *RequestAddKey) Validate() error {
	if r.Name == "" {
		return errors.New("name cannot be empty")
	}
	if r.Mnemonic == "" {
		return errors.New("mnemonic cannot be empty")
	}
	if !bip39.IsMnemonicValid(r.Mnemonic) {
		return fmt.Errorf("invalid mnemonic %s", r.Mnemonic)
	}

	return nil
}
