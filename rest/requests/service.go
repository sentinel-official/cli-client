package requests

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/pkg/errors"
	hubtypes "github.com/sentinel-official/hub/types"
)

type Connect struct {
	Backend  string `json:"backend"`
	Password string `json:"password"`

	ID   uint64 `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`

	Info      []byte   `json:"info"`
	Keys      [][]byte `json:"keys"`
	Resolvers []net.IP `json:"resolvers"`
}

func NewConnect(r *http.Request) (*Connect, error) {
	var v Connect
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (r *Connect) Validate() error {
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

	if r.ID == 0 {
		return errors.New("id cannot be 0")
	}
	if r.From == "" {
		return errors.New("from cannot be empty")
	}
	if r.To == "" {
		return errors.New("to cannot be empty")
	}
	if _, err := hubtypes.NodeAddressFromBech32(r.To); err != nil {
		return errors.Wrap(err, "invalid to")
	}

	if r.Info == nil {
		return errors.New("info cannot be nil")
	}
	if len(r.Info) != 58 {
		return errors.New("info length must be 58 bytes")
	}
	if r.Keys == nil {
		return errors.New("keys cannot be nil")
	}
	if len(r.Keys) != 1 {
		return errors.New("keys length must be 1")
	}
	if len(r.Keys[0]) != 32 {
		return errors.New("key at index 0 length must be 32 bytes")
	}

	return nil
}
