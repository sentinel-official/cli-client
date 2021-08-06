package requests

import (
	"encoding/json"
	"fmt"
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
		return fmt.Errorf("backend must be one of [%s, %s, %s]",
			keyring.BackendFile, keyring.BackendOS, keyring.BackendTest)
	}
	if r.Backend == keyring.BackendFile {
		if r.Password == "" {
			return errors.New("password cannot be empty")
		}
		if len(r.Password) < 8 {
			return fmt.Errorf("password length cannot be less than %d", 8)
		}
	}

	if r.ID == 0 {
		return errors.New("id cannot be zero")
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
	if len(r.Info) != 4+16+4+2+32 {
		return fmt.Errorf("info length must be %d bytes", 4+16+4+2+32)
	}
	if r.Keys == nil {
		return errors.New("keys cannot be nil")
	}
	if len(r.Keys) != 1 {
		return fmt.Errorf("keys length must be %d", 1)
	}
	if len(r.Keys[0]) != 32 {
		return fmt.Errorf("key at index %d length must be %d bytes", 0, 32)
	}

	return nil
}
