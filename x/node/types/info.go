package types

import (
	"time"

	clienttypes "github.com/sentinel-official/cli-client/types"
)

type (
	Handshake struct {
		Enable bool   `json:"enable"`
		Peers  uint64 `json:"peers"`
	}
	Location struct {
		City      string  `json:"city"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
)

type (
	Info struct {
		Address                string                `json:"address"`
		Bandwidth              clienttypes.Bandwidth `json:"bandwidth"`
		Handshake              Handshake             `json:"handshake"`
		IntervalSetSessions    time.Duration         `json:"interval_set_sessions"`
		IntervalUpdateSessions time.Duration         `json:"interval_update_sessions"`
		IntervalUpdateStatus   time.Duration         `json:"interval_update_status"`
		Location               Location              `json:"location"`
		Moniker                string                `json:"moniker"`
		Operator               string                `json:"operator"`
		Peers                  int                   `json:"peers"`
		Price                  string                `json:"price"`
		Provider               string                `json:"provider"`
		Type                   uint64                `json:"type"`
		Version                string                `json:"version"`
	}
)
