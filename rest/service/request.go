package service

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/pkg/errors"
	hubtypes "github.com/sentinel-official/hub/types"
)

type RequestConnect struct {
	ID   uint64 `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`

	Info      []byte   `json:"info"`
	Keys      [][]byte `json:"keys"`
	Resolvers []net.IP `json:"resolvers"`
}

func NewRequestConnect(r *http.Request) (*RequestConnect, error) {
	var v RequestConnect
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (r *RequestConnect) Validate() error {
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
