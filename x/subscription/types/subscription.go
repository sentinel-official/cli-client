package types

import (
	"time"

	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"
)

type Subscription struct {
	ID       uint64    `json:"id"`
	Address  string    `json:"address"`
	ExpiryAt time.Time `json:"expiry_at"`
	Status   string    `json:"status"`
	StatusAt time.Time `json:"status_at"`

	NodeAddress string `json:"node_address"`
	Gigabytes   int64  `json:"gigabytes"`
	Hours       int64  `json:"hours"`
	Deposit     string `json:"deposit"`

	PlanID uint64 `json:"plan_id"`
	Denom  string `json:"denom"`
}

func NewSubscriptionFromRaw(v subscriptiontypes.Subscription) Subscription {
	s := Subscription{
		ID:       v.GetID(),
		Address:  v.GetAddress().String(),
		ExpiryAt: v.GetExpiryAt(),
		Status:   v.GetStatus().String(),
		StatusAt: v.GetStatusAt(),
	}

	if v.Type() == subscriptiontypes.TypeNode {
		s.NodeAddress = v.(*subscriptiontypes.NodeSubscription).NodeAddress
		s.Gigabytes = v.(*subscriptiontypes.NodeSubscription).Gigabytes
		s.Hours = v.(*subscriptiontypes.NodeSubscription).Hours
		s.Deposit = v.(*subscriptiontypes.NodeSubscription).Deposit.String()
	}
	if v.Type() == subscriptiontypes.TypePlan {
		s.PlanID = v.(*subscriptiontypes.PlanSubscription).PlanID
		s.Denom = v.(*subscriptiontypes.PlanSubscription).Denom
	}

	return s
}

type Subscriptions []Subscription

func NewSubscriptionsFromRaw(v subscriptiontypes.Subscriptions) Subscriptions {
	items := make(Subscriptions, 0, len(v))
	for i := 0; i < len(v); i++ {
		items = append(items, NewSubscriptionFromRaw(v[i]))
	}

	return items
}
