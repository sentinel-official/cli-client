package cmd

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	hubtypes "github.com/sentinel-official/hub/types"
	nodetypes "github.com/sentinel-official/hub/x/node/types"
	"github.com/spf13/cobra"
)

func GetTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Node related subcommands",
	}

	cmd.AddCommand(
		txSetStatus(),
	)

	return cmd
}

func txSetStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status-set [status]",
		Short: "Set a node status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := nodetypes.NewMsgSetStatusRequest(
				ctx.FromAddress.Bytes(),
				hubtypes.StatusFromString(args[0]),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
