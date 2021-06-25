package cmd

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	nodecmd "github.com/sentinel-official/cli-client/x/node/client/cmd"
	plancmd "github.com/sentinel-official/cli-client/x/plan/client/cmd"
	providercmd "github.com/sentinel-official/cli-client/x/provider/client/cmd"
	sessioncmd "github.com/sentinel-official/cli-client/x/session/client/cmd"
	subscriptioncmd "github.com/sentinel-official/cli-client/x/subscription/client/cmd"
)

func TxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(nodecmd.GetTxCommand())
	cmd.AddCommand(plancmd.GetTxCommand())
	cmd.AddCommand(providercmd.GetTxCommand())
	cmd.AddCommand(subscriptioncmd.GetTxCommand())
	cmd.AddCommand(sessioncmd.GetTxCommand())

	cmd.PersistentFlags().String(flags.FlagChainID, "", "the network chain identity")

	return cmd
}
