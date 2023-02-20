package types

import (
	"encoding/json"
	"os"
)

type Service interface {
	Info() []byte
	PreUp() error
	IsUp() bool
	Up() error
	PostUp() error
	PreDown() error
	Down() error
	PostDown() error
	Transfer() (int64, int64, error)
}

type Status struct {
	From string `json:"from"`
	ID   uint64 `json:"id"`
	To   string `json:"to"`
	Type uint64 `json:"type"`
	Info []byte `json:"info"`
}

func NewStatus() *Status {
	return &Status{}
}

func (s *Status) WithFrom(v string) *Status { s.From = v; return s }
func (s *Status) WithID(v uint64) *Status   { s.ID = v; return s }
func (s *Status) WithInfo(v []byte) *Status { s.Info = v; return s }
func (s *Status) WithTo(v string) *Status   { s.To = v; return s }
func (s *Status) WithType(v uint64) *Status { s.Type = v; return s }

func (s *Status) LoadFromPath(path string) error {
	if _, err := os.Stat(path); err != nil {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, s)
}

func (s *Status) SaveToPath(path string) error {
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0600)
}
