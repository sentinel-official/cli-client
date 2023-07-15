package types

import (
	"time"

	plantypes "github.com/sentinel-official/hub/x/plan/types"

	clienttypes "github.com/sentinel-official/cli-client/types"
)

type Plan struct {
	ID       uint64            `json:"id"`
	Address  string            `json:"address"`
	Prices   clienttypes.Coins `json:"prices"`
	Duration time.Duration     `json:"duration"`
	Bytes    int64             `json:"bytes"`
	Status   string            `json:"status"`
	StatusAt time.Time         `json:"status_at"`
}

func NewPlanFromRaw(v *plantypes.Plan) Plan {
	return Plan{
		ID:       v.ID,
		Address:  v.Address,
		Prices:   clienttypes.NewCoinsFromRaw(v.Prices),
		Duration: v.Duration,
		Bytes:    v.Bytes.Int64(),
		Status:   v.Status.String(),
		StatusAt: v.StatusAt,
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
