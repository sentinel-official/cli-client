package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/spf13/cobra"
)

const (
	FlagAccount  = "account"
	FlagCoinType = "coin-type"
	FlagIndex    = "index"
	FlagRating   = "rating"
	FlagRecover  = "recover"
	FlagResolver = "resolver"
)

const (
	FlagHome        = "home"
	FlagListen      = "listen"
	FlagTTY         = "tty"
	FlagWithKeyring = "with-keyring"
	FlagWithService = "with-service"
)

const (
	FlagBroadcastMode  = "broadcast-mode"
	FlagChainID        = "chain-id"
	FlagFrom           = "from"
	FlagGas            = "gas"
	FlagGasPrices      = "gas-prices"
	FlagKeyringBackend = "keyring-backend"
	FlagKeyringHome    = "keyring-home"
	FlagMemo           = "memo"
	FlagServiceHome    = "service-home"
	FlagRPCAddress     = "rpc-address"
	FlagTimeout        = "timeout"
)

func AddFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagBroadcastMode, "block", "transaction broadcasting mode (sync|async|block)")
	cmd.Flags().String(FlagChainID, "", "chain identity of the network")
	cmd.Flags().String(FlagFrom, "", "name or address of private key with which to sign")
	cmd.Flags().Uint64(FlagGas, 200000, "gas limit to set per-transaction")
	cmd.Flags().String(FlagGasPrices, "", "gas prices in decimal format to determine the transaction fee")
	cmd.Flags().String(FlagKeyringBackend, keyring.BackendOS, "the keyring backend backend (os|file|test)")
	cmd.Flags().String(FlagKeyringHome, Home, "home directory of the keys")
	cmd.Flags().String(FlagMemo, "", "memo to send along with transaction")
	cmd.Flags().String(FlagServiceHome, Home, "home directory of the service")
	cmd.Flags().String(FlagRPCAddress, "", "tendermint RPC interface address for this chain")
	cmd.Flags().Duration(FlagTimeout, 15*time.Second, "time limit for requests made by the HTTP client")
}

func AddKeyringFlagsToCmd(cmd *cobra.Command) {
	AddFlagsToCmd(cmd)
	_ = cmd.Flags().MarkHidden(FlagBroadcastMode)
	_ = cmd.Flags().MarkHidden(FlagChainID)
	_ = cmd.Flags().MarkHidden(FlagFrom)
	_ = cmd.Flags().MarkHidden(FlagGas)
	_ = cmd.Flags().MarkHidden(FlagGasPrices)
	_ = cmd.Flags().MarkHidden(FlagMemo)
	_ = cmd.Flags().MarkHidden(FlagServiceHome)
	_ = cmd.Flags().MarkHidden(FlagRPCAddress)
}

func AddQueryFlagsToCmd(cmd *cobra.Command) {
	AddFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(FlagRPCAddress)
	_ = cmd.Flags().MarkHidden(FlagBroadcastMode)
	_ = cmd.Flags().MarkHidden(FlagChainID)
	_ = cmd.Flags().MarkHidden(FlagFrom)
	_ = cmd.Flags().MarkHidden(FlagGas)
	_ = cmd.Flags().MarkHidden(FlagGasPrices)
	_ = cmd.Flags().MarkHidden(FlagKeyringBackend)
	_ = cmd.Flags().MarkHidden(FlagKeyringHome)
	_ = cmd.Flags().MarkHidden(FlagMemo)
	_ = cmd.Flags().MarkHidden(FlagServiceHome)
}

func AddTxFlagsToCmd(cmd *cobra.Command) {
	AddFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(FlagChainID)
	_ = cmd.MarkFlagRequired(FlagFrom)
	_ = cmd.MarkFlagRequired(FlagRPCAddress)
}
