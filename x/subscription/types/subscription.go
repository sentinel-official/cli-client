package types

import (
	"time"

	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"

	clienttypes "github.com/sentinel-official/cli-client/types"
)

type Subscription struct {
	ID       uint64           `json:"id"`
	Owner    string           `json:"owner"`
	Plan     uint64           `json:"plan"`
	Expiry   time.Time        `json:"expiry"`
	Denom    string           `json:"denom"`
	Node     string           `json:"node"`
	Price    clienttypes.Coin `json:"price"`
	Deposit  clienttypes.Coin `json:"deposit"`
	Free     int64            `json:"free"`
	Status   string           `json:"status"`
	StatusAt time.Time        `json:"status_at"`
}

func NewSubscriptionFromRaw(v *subscriptiontypes.Subscription) Subscription {
	return Subscription{
		ID:       v.Id,
		Owner:    v.Owner,
		Plan:     v.Plan,
		Expiry:   v.Expiry,
		Denom:    v.Denom,
		Node:     v.Node,
		Price:    clienttypes.NewCoinFromRaw(&v.Price),
		Deposit:  clienttypes.NewCoinFromRaw(&v.Deposit),
		Free:     v.Free.Int64(),
		Status:   v.Status.String(),
		StatusAt: v.StatusAt,
	}
}

type Subscriptions []Subscription

func NewSubscriptionsFromRaw(v subscriptiontypes.Subscriptions) Subscriptions {
	items := make(Subscriptions, 0, len(v))
	for i := 0; i < len(v); i++ {
		items = append(items, NewSubscriptionFromRaw(&v[i]))
	}

	return items
}
