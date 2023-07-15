package types

import (
	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"
)

type Allocation struct {
	Address       string `json:"address"`
	GrantedBytes  int64  `json:"granted_bytes"`
	UtilisedBytes int64  `json:"utilised_bytes"`
}

func NewAllocationFromRaw(v *subscriptiontypes.Allocation) Allocation {
	return Allocation{
		Address:       v.Address,
		GrantedBytes:  v.GrantedBytes.Int64(),
		UtilisedBytes: v.UtilisedBytes.Int64(),
	}
}

type Allocations []Allocation

func NewAllocationsFromRaw(v subscriptiontypes.Allocations) Allocations {
	items := make(Allocations, 0, len(v))
	for i := 0; i < len(v); i++ {
		items = append(items, NewAllocationFromRaw(&v[i]))
	}

	return items
}
