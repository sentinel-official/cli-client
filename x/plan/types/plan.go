package types

import (
	"time"

	plantypes "github.com/sentinel-official/hub/x/plan/types"

	clienttypes "github.com/sentinel-official/cli-client/types"
)

type Plan struct {
	ID       uint64            `json:"id"`
	Provider string            `json:"provider"`
	Price    clienttypes.Coins `json:"price"`
	Validity time.Duration     `json:"validity"`
	Bytes    int64             `json:"bytes"`
	Status   string            `json:"status"`
	StatusAt time.Time         `json:"status_at"`
}

func NewPlanFromRaw(v *plantypes.Plan) Plan {
	return Plan{
		ID:       v.Id,
		Provider: v.Provider,
		Price:    clienttypes.NewCoinsFromRaw(v.Price),
		Validity: v.Validity,
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
