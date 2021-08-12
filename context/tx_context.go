package context

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	clitypes "github.com/sentinel-official/cli-client/types"
)

type TxContext struct {
	QueryContext
	Gas       uint64
	GasPrices sdk.DecCoins
	Memo      string
}

func NewTxContextFromCmd(cmd *cobra.Command) (ctx TxContext, err error) {
	ctx.QueryContext, err = NewQueryContextFromCmd(cmd)
	if err != nil {
		return ctx, err
	}

	ctx.BroadcastMode, err = cmd.Flags().GetString(clitypes.FlagBroadcastMode)
	if err != nil {
		return ctx, err
	}

	ctx.ChainID, err = cmd.Flags().GetString(clitypes.FlagChainID)
	if err != nil {
		return ctx, err
	}

	ctx.From, err = cmd.Flags().GetString(clitypes.FlagFrom)
	if err != nil {
		return ctx, err
	}

	ctx.Gas, err = cmd.Flags().GetUint64(clitypes.FlagGas)
	if err != nil {
		return ctx, err
	}

	s, err := cmd.Flags().GetString(clitypes.FlagGasPrices)
	if err != nil {
		return ctx, err
	}

	ctx.GasPrices, err = sdk.ParseDecCoins(s)
	if err != nil {
		return ctx, err
	}

	ctx.Memo, err = cmd.Flags().GetString(clitypes.FlagMemo)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (c TxContext) WithQueryContext(v QueryContext) TxContext {
	c.QueryContext = v
	return c
}

func (c TxContext) WithGas(v uint64) TxContext {
	c.Gas = v
	return c
}

func (c TxContext) WithGasPrices(v sdk.DecCoins) TxContext {
	c.GasPrices = v
	return c
}

func (c TxContext) WithMemo(v string) TxContext {
	c.Memo = v
	return c
}
