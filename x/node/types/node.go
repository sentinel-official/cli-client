package types

import (
	"time"

	nodetypes "github.com/sentinel-official/hub/x/node/types"

	clienttypes "github.com/sentinel-official/cli-client/types"
)

type Node struct {
	Info
	Address        string            `json:"address"`
	GigabytePrices clienttypes.Coins `json:"gigabyte_prices"`
	HourlyPrices   clienttypes.Coins `json:"hourly_prices"`
	RemoteURL      string            `json:"remote_url"`
	Status         string            `json:"status"`
	StatusAt       time.Time         `json:"status_at"`
}

func (n Node) WithInfo(v Info) Node { n.Info = v; return n }

func NewNodeFromRaw(v *nodetypes.Node) Node {
	return Node{
		Address:        v.Address,
		GigabytePrices: clienttypes.NewCoinsFromRaw(v.GigabytePrices),
		HourlyPrices:   clienttypes.NewCoinsFromRaw(v.HourlyPrices),
		RemoteURL:      v.RemoteURL,
		Status:         v.Status.String(),
		StatusAt:       v.StatusAt,
	}
}

type Nodes []Node

func NewNodesFromRaw(v nodetypes.Nodes) Nodes {
	items := make(Nodes, 0, len(v))
	for i := 0; i < len(v); i++ {
		items = append(items, NewNodeFromRaw(&v[i]))
	}

	return items
}
