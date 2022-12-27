package types

import (
	"time"

	nodetypes "github.com/sentinel-official/hub/x/node/types"

	clitypes "github.com/sentinel-official/cli-client/types"
)

type Node struct {
	NodeInfo
	Address   string         `json:"address"`
	Provider  string         `json:"provider"`
	Price     clitypes.Coins `json:"price"`
	RemoteURL string         `json:"remote_url"`
	Status    string         `json:"status"`
	StatusAt  time.Time      `json:"status_at"`
}

func (n Node) WithInfo(v NodeInfo) Node { n.NodeInfo = v; return n }

func NewNodeFromRaw(v *nodetypes.Node) Node {
	return Node{
		Address:   v.Address,
		Provider:  v.Provider,
		Price:     clitypes.NewCoinsFromRaw(v.Price),
		RemoteURL: v.RemoteURL,
		Status:    v.Status.String(),
		StatusAt:  v.StatusAt,
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
