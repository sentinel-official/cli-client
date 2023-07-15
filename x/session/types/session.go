package types

import (
	"time"

	sessiontypes "github.com/sentinel-official/hub/x/session/types"

	clienttypes "github.com/sentinel-official/cli-client/types"
)

type Session struct {
	ID             uint64                `json:"id"`
	SubscriptionID uint64                `json:"subscription_id"`
	NodeAddress    string                `json:"node_address"`
	Address        string                `json:"address"`
	Duration       time.Duration         `json:"duration"`
	Bandwidth      clienttypes.Bandwidth `json:"bandwidth"`
	Status         string                `json:"status"`
	StatusAt       time.Time             `json:"status_at"`
}

func NewSessionFromRaw(v *sessiontypes.Session) Session {
	return Session{
		ID:             v.ID,
		SubscriptionID: v.SubscriptionID,
		NodeAddress:    v.NodeAddress,
		Address:        v.Address,
		Duration:       v.Duration,
		Bandwidth:      clienttypes.NewBandwidthFromRaw(v.Bandwidth),
		Status:         v.Status.String(),
		StatusAt:       v.StatusAt,
	}
}

type Sessions []Session

func NewSessionsFromRaw(v sessiontypes.Sessions) Sessions {
	items := make(Sessions, 0, len(v))
	for i := 0; i < len(v); i++ {
		items = append(items, NewSessionFromRaw(&v[i]))
	}

	return items
}
