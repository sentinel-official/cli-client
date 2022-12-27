package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	hubtypes "github.com/sentinel-official/hub/types"
	"github.com/spf13/cobra"
)

const (
	FlagAccount        = "account"
	FlagBroadcastMode  = "broadcast-mode"
	FlagChainID        = "chain-id"
	FlagCoinType       = "coin-type"
	FlagDescription    = "description"
	FlagFrom           = "from"
	FlagGas            = "gas"
	FlagGasPrices      = "gas-prices"
	FlagHome           = "home"
	FlagIdentity       = "identity"
	FlagIndex          = "index"
	FlagKeyringBackend = "keyring.backend"
	FlagKeyringHome    = "keyring.home"
	FlagListen         = "listen"
	FlagMemo           = "memo"
	FlagName           = "name"
	FlagProvider       = "provider"
	FlagRating         = "rating"
	FlagRecover        = "recover"
	FlagResolver       = "resolver"
	FlagRPCAddress     = "rpc-address"
	FlagServiceHome    = "service.home"
	FlagStatus         = "status"
	FlagTimeout        = "timeout"
	FlagTTY            = "tty"
	FlagWebsite        = "website"
	FlagWithKeyring    = "with-keyring"
	FlagWithService    = "with-service"
	FlagAddress        = keys.FlagAddress
)

func addKeyringFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagKeyringBackend, keyring.BackendOS, "the keyring backend (file|os|test)")
	cmd.Flags().String(FlagKeyringHome, Home, "home directory of the keyring")
}

func addQueryFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagRPCAddress, "", "tendermint RPC interface address for this chain")
	_ = cmd.MarkFlagRequired(FlagRPCAddress)
}

func addServiceFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagServiceHome, Home, "home directory of the service")
}

func addTimeoutFlagsToCmd(cmd *cobra.Command) {
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
	addTimeoutFlagsToCmd(cmd)
}

func AddPaginationFlagsToCmd(cmd *cobra.Command, query string) {
	flags.AddPaginationFlagsToCmd(cmd, query)
}

func AddQueryFlagsToCmd(cmd *cobra.Command) {
	addQueryFlagsToCmd(cmd)
}

func AddServiceFlagsToCmd(cmd *cobra.Command) {
	addServiceFlagsToCmd(cmd)
}

func AddTimeoutFlagsToCmd(cmd *cobra.Command) {
	addTimeoutFlagsToCmd(cmd)
}

func AddTxFlagsToCmd(cmd *cobra.Command) {
	addKeyringFlagsToCmd(cmd)
	addQueryFlagsToCmd(cmd)
	addTimeoutFlagsToCmd(cmd)
	addTxFlagsToCmd(cmd)
}

func GetAccAddressFromCmd(cmd *cobra.Command) (sdk.AccAddress, error) {
	s, err := cmd.Flags().GetString(FlagAddress)
	if err != nil {
		return nil, err
	}
	if s == "" {
		return nil, nil
	}

	return sdk.AccAddressFromBech32(s)
}

func GetPageRequestFromCmd(cmd *cobra.Command) (*query.PageRequest, error) {
	return client.ReadPageRequest(cmd.Flags())
}

func GetProvAddressFromCmd(cmd *cobra.Command) (hubtypes.ProvAddress, error) {
	s, err := cmd.Flags().GetString(FlagProvider)
	if err != nil {
		return nil, err
	}
	if s == "" {
		return nil, nil
	}

	return hubtypes.ProvAddressFromBech32(s)
}

func GetStatusFromCmd(cmd *cobra.Command) (hubtypes.Status, error) {
	s, err := cmd.Flags().GetString(FlagStatus)
	if err != nil {
		return 0, err
	}

	return hubtypes.StatusFromString(s), nil
}
