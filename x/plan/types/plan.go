package types

import (
	"time"

	plantypes "github.com/sentinel-official/hub/x/plan/types"

	clienttypes "github.com/sentinel-official/cli-client/types"
)

type Plan struct {
	ID        uint64            `json:"id"`
	Address   string            `json:"address"`
	Prices    clienttypes.Coins `json:"prices"`
	Duration  time.Duration     `json:"duration"`
	Gigabytes int64             `json:"gigabytes"`
	Status    string            `json:"status"`
	StatusAt  time.Time         `json:"status_at"`
}

func NewPlanFromRaw(v *plantypes.Plan) Plan {
	return Plan{
		ID:        v.ID,
		Address:   v.ProviderAddress,
		Prices:    clienttypes.NewCoinsFromRaw(v.Prices),
		Duration:  v.Duration,
		Gigabytes: v.Gigabytes,
		Status:    v.Status.String(),
		StatusAt:  v.StatusAt,
	}
}

type Plans []Plan

func NewPlansFromRaw(v plantypes.Plans) Plans {
	items := make(Plans, 0, len(v))
	for i := 0; i < len(v); i++ {
		items = append(items, NewPlanFromRaw(&v[i]))
	}

	return items
}
