package context

import (
	"encoding/base64"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/spf13/cobra"

	clitypes "github.com/sentinel-official/cli-client/types"
)

type TxContext struct {
	KeyringContext
	QueryContext
	Gas       uint64
	GasPrices sdk.DecCoins
	Memo      string
}

func NewTxContextFromCmd(cmd *cobra.Command) (ctx TxContext, err error) {
	ctx.KeyringContext, err = NewKeyringContextFromCmd(cmd)
	if err != nil {
		return ctx, err
	}

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

func (c *TxContext) SignMessagesAndBroadcastTx(password string, messages ...sdk.Msg) (*sdk.TxResponse, error) {
	key, err := c.GetKey(password, c.From)
	if err != nil {
		return nil, err
	}

	accAddr, err := base64.StdEncoding.DecodeString(key.Address)
	if err != nil {
		return nil, err
	}

	pubKey, err := base64.StdEncoding.DecodeString(key.PubKey)
	if err != nil {
		return nil, err
	}

	account, err := c.QueryAccount(accAddr)
	if err != nil {
		return nil, err
	}

	txb := c.TxConfig.NewTxBuilder()
	if err := txb.SetMsgs(messages...); err != nil {
		return nil, err
	}

	txb.SetGasLimit(c.Gas)
	txb.SetMemo(c.Memo)

	if !c.GasPrices.IsZero() {
		var (
			gas  = sdk.NewDec(int64(c.Gas))
			fees = make(sdk.Coins, len(c.GasPrices))
		)

		for i, price := range c.GasPrices {
			fee := price.Amount.Mul(gas)
			fees[i] = sdk.NewCoin(price.Denom, fee.Ceil().RoundInt())
		}

		txb.SetFeeAmount(fees)
	}

	txSignature := txsigning.SignatureV2{
		PubKey: &secp256k1.PubKey{
			Key: pubKey,
		},
		Data: &txsigning.SingleSignatureData{
			SignMode:  c.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
		Sequence: account.GetSequence(),
	}

	if err := txb.SetSignatures(txSignature); err != nil {
		return nil, err
	}

	message, err := c.TxConfig.SignModeHandler().GetSignBytes(
		c.TxConfig.SignModeHandler().DefaultMode(),
		authsigning.SignerData{
			ChainID:       c.ChainID,
			AccountNumber: account.GetAccountNumber(),
			Sequence:      account.GetSequence(),
		},
		txb.GetTx(),
	)
	if err != nil {
		return nil, err
	}

	res, err := c.SignMessage(password, c.From, message)
	if err != nil {
		return nil, err
	}

	signature, err := base64.StdEncoding.DecodeString(res.Signature)
	if err != nil {
		return nil, err
	}

	txSignature.Data = &txsigning.SingleSignatureData{
		SignMode:  c.TxConfig.SignModeHandler().DefaultMode(),
		Signature: signature,
	}

	if err := txb.SetSignatures(txSignature); err != nil {
		return nil, err
	}

	txBytes, err := c.TxConfig.TxEncoder()(txb.GetTx())
	if err != nil {
		return nil, err
	}

	return c.QueryContext.BroadcastTx(txBytes)
}
