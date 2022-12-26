package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
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
	FlagKeyringBackend = "keyring.backend"
	FlagKeyringHome    = "keyring.home"
	FlagMemo           = "memo"
	FlagServiceHome    = "service.home"
	FlagRPCAddress     = "rpc-address"
	FlagTimeout        = "timeout"
)

func addKeyringFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagKeyringBackend, keyring.BackendOS, "the keyring backend (file|os|test)")
	cmd.Flags().String(FlagKeyringHome, Home, "home directory of the keyring")
}

func addQueryFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagRPCAddress, "", "tendermint RPC interface address for this chain")
	_ = cmd.MarkFlagRequired(FlagRPCAddress)
}

func addTimeoutFlagToCmd(cmd *cobra.Command) {
	cmd.Flags().Duration(FlagTimeout, 15*time.Second, "time limit for requests made by the HTTP client")
}

func addTxFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagBroadcastMode, flags.BroadcastBlock, "transaction broadcasting mode (async|block|sync)")
	cmd.Flags().String(FlagChainID, "", "chain identity of the network")
	cmd.Flags().String(FlagFrom, "", "name or address of private key with which to sign")
	cmd.Flags().Uint64(FlagGas, flags.DefaultGasLimit, "gas limit to set per-transaction")
	cmd.Flags().String(FlagGasPrices, "", "gas prices in decimal format to determine the transaction fee")
	cmd.Flags().String(FlagMemo, "", "memo to send along with transaction")
	_ = cmd.MarkFlagRequired(FlagChainID)
	_ = cmd.MarkFlagRequired(FlagFrom)
}

func AddKeyringFlagsToCmd(cmd *cobra.Command) {
	addKeyringFlagsToCmd(cmd)
	addTimeoutFlagToCmd(cmd)
}

func AddQueryFlagsToCmd(cmd *cobra.Command) {
	addQueryFlagsToCmd(cmd)
	addTimeoutFlagToCmd(cmd)
}

func AddServiceFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagServiceHome, Home, "home directory of the service")
}

func AddTxFlagsToCmd(cmd *cobra.Command) {
	addKeyringFlagsToCmd(cmd)
	addQueryFlagsToCmd(cmd)
	addTimeoutFlagToCmd(cmd)
	addTxFlagsToCmd(cmd)
}
