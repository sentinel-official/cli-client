package requests

import (
	"encoding/json"
	"net/http"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/go-bip39"
	"github.com/pkg/errors"
)

type GeyKey struct {
	Backend  string `json:"backend"`
	Password string `json:"password"`

	Name string `json:"name"`
}

func NewGeyKey(r *http.Request) (*GeyKey, error) {
	var v GeyKey
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (r *GeyKey) Validate() error {
	if r.Backend == "" {
		return errors.New("backend cannot be empty")
	}
	if r.Backend != keyring.BackendFile && r.Backend != keyring.BackendOS && r.Backend != keyring.BackendTest {
		return errors.New("backend must be either file, os, or test")
	}
	if r.Backend == keyring.BackendFile {
		if r.Password == "" {
			return errors.New("password cannot be empty")
		}
		if len(r.Password) < 8 {
			return errors.New("password length cannot be less than 8 characters")
		}
	}

	if r.Name == "" {
		return errors.New("name cannot be empty")
	}

	return nil
}

type GeyKeys struct {
	Backend  string ` json:"backend"`
	Password string `json:"password"`
}

func NewGeyKeys(r *http.Request) (*GeyKeys, error) {
	var v GeyKeys
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (r *GeyKeys) Validate() error {
	if r.Backend == "" {
		return errors.New("backend cannot be empty")
	}
	if r.Backend != keyring.BackendFile && r.Backend != keyring.BackendOS && r.Backend != keyring.BackendTest {
		return errors.New("backend must be either file, os, or test")
	}
	if r.Backend == keyring.BackendFile {
		if r.Password == "" {
			return errors.New("password cannot be empty")
		}
		if len(r.Password) < 8 {
			return errors.New("password length cannot be less than 8 characters")
		}
	}

	return nil
}

type SignMessage struct {
	Backend  string `json:"backend"`
	Password string `json:"password"`

	Name    string `json:"name"`
	Message []byte `json:"message"`
}

func NewSignMessage(r *http.Request) (*SignMessage, error) {
	var v SignMessage
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (r *SignMessage) Validate() error {
	if r.Backend == "" {
		return errors.New("backend cannot be empty")
	}
	if r.Backend != keyring.BackendFile && r.Backend != keyring.BackendOS && r.Backend != keyring.BackendTest {
		return errors.New("backend must be either file, os, or test")
	}
	if r.Backend == keyring.BackendFile {
		if r.Password == "" {
			return errors.New("password cannot be empty")
		}
		if len(r.Password) < 8 {
			return errors.New("password length cannot be less than 8 characters")
		}
	}

	if r.Name == "" {
		return errors.New("name cannot be empty")
	}

	return nil
}

type AddKey struct {
	Backend  string `json:"backend"`
	Password string `json:"password"`

	Name          string `json:"name"`
	Mnemonic      string `json:"mnemonic"`
	CoinType      uint32 `json:"coin_type"`
	Account       uint32 `json:"account"`
	Index         uint32 `json:"index"`
	BIP39Password string `json:"bip39_password"`
}

func NewAddKey(r *http.Request) (*AddKey, error) {
	var v AddKey
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (r *AddKey) Validate() error {
	if r.Backend == "" {
		return errors.New("backend cannot be empty")
	}
	if r.Backend != keyring.BackendFile && r.Backend != keyring.BackendOS && r.Backend != keyring.BackendTest {
		return errors.New("backend must be either file, os, or test")
	}
	if r.Backend == keyring.BackendFile {
		if r.Password == "" {
			return errors.New("password cannot be empty")
		}
		if len(r.Password) < 8 {
			return errors.New("password length cannot be less than 8 characters")
		}
	}

	if r.Name == "" {
		return errors.New("name cannot be empty")
	}
	if r.Mnemonic == "" {
		return errors.New("mnemonic cannot be empty")
	}
	if !bip39.IsMnemonicValid(r.Mnemonic) {
		return errors.New("invalid mnemonic")
	}

	return nil
}

type DeleteKey struct {
	Backend  string `json:"backend"`
	Password string `json:"password"`

	Name string `json:"name"`
}

func NewDeleteKey(r *http.Request) (*DeleteKey, error) {
	var v DeleteKey
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (r *DeleteKey) Validate() error {
	if r.Backend == "" {
		return errors.New("backend cannot be empty")
	}
	if r.Backend != keyring.BackendFile && r.Backend != keyring.BackendOS && r.Backend != keyring.BackendTest {
		return errors.New("backend must be either file, os, or test")
	}
	if r.Backend == keyring.BackendFile {
		if r.Password == "" {
			return errors.New("password cannot be empty")
		}
		if len(r.Password) < 8 {
			return errors.New("password length cannot be less than 8 characters")
		}
	}

	if r.Name == "" {
		return errors.New("name cannot be empty")
	}

	return nil
}
