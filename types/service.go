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

type ServiceStatus struct {
	ID    uint64 `json:"id"`
	IFace string `json:"iface"`
}

func NewServiceStatus() *ServiceStatus {
	return &ServiceStatus{}
}

func (s *ServiceStatus) WithID(v uint64) *ServiceStatus    { s.ID = v; return s }
func (s *ServiceStatus) WithIFace(v string) *ServiceStatus { s.IFace = v; return s }

func (s *ServiceStatus) LoadFromPath(path string) error {
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

func (s *ServiceStatus) SaveToPath(path string) error {
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0600)
}
