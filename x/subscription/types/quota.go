package types

import (
	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"
)

type Quota struct {
	Address   string `json:"address"`
	Consumed  int64  `json:"consumed"`
	Allocated int64  `json:"allocated"`
}

func NewQuotaFromRaw(v *subscriptiontypes.Quota) Quota {
	return Quota{
		Address:   v.Address,
		Consumed:  v.Consumed.Int64(),
		Allocated: v.Allocated.Int64(),
	}
}

type Quotas []Quota

func NewQuotasFromRaw(v subscriptiontypes.Quotas) Quotas {
	items := make(Quotas, 0, len(v))
	for i := 0; i < len(v); i++ {
		items = append(items, NewQuotaFromRaw(&v[i]))
	}

	return items
}
