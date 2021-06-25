package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Coin struct {
	Denom string `json:"denom"`
	Value int64  `json:"value"`
}

func NewCoinFromRaw(v *sdk.Coin) Coin {
	return Coin{
		Denom: v.Denom,
		Value: v.Amount.Int64(),
	}
}

func (c *Coin) Raw() sdk.Coin {
	return sdk.Coin{
		Denom:  c.Denom,
		Amount: sdk.NewInt(c.Value),
	}
}

type Coins []Coin

func NewCoinsFromRaw(v sdk.Coins) Coins {
	items := make(Coins, 0, v.Len())
	for i := 0; i < len(v); i++ {
		items = append(items, NewCoinFromRaw(&v[i]))
	}

	return items
}

func (c Coins) Raw() sdk.Coins {
	items := make(sdk.Coins, 0, len(c))
	for i := 0; i < len(c); i++ {
		items = append(items, c[i].Raw())
	}

	return items
}

type DecCoin struct {
	Denom string `json:"denom"`
	Value string `json:"value"`
}

func NewDecCoinFromRaw(v *sdk.DecCoin) DecCoin {
	return DecCoin{
		Denom: v.Denom,
		Value: v.Amount.String(),
	}
}

func (c *DecCoin) Raw() sdk.DecCoin {
	return sdk.DecCoin{
		Denom:  c.Denom,
		Amount: sdk.MustNewDecFromStr(c.Denom),
	}
}

type DecCoins []DecCoin

func NewDecCoinsFromRaw(v sdk.DecCoins) DecCoins {
	items := make(DecCoins, 0, v.Len())
	for i := 0; i < len(v); i++ {
		items = append(items, NewDecCoinFromRaw(&v[i]))
	}

	return items
}

func (c DecCoins) Raw() sdk.DecCoins {
	items := make(sdk.DecCoins, 0, len(c))
	for i := 0; i < len(c); i++ {
		items = append(items, c[i].Raw())
	}

	return items
}
