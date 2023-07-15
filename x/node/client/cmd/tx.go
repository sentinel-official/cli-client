package cmd

import (
	"strconv"

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
		txUpdateStatus(),
		txSubscribe(),
	)

	return cmd
}

func txUpdateStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-status [status]",
		Short: "Update a node status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := nodetypes.NewMsgUpdateStatusRequest(
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

func txSubscribe() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscribe [node-addr] [gigabytes] [hours] [denom]",
		Short: "Subscribe to a node",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			address, err := hubtypes.NodeAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			gigabytes, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			hours, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := nodetypes.NewMsgSubscribeRequest(
				ctx.FromAddress,
				address,
				gigabytes,
				hours,
				args[3],
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
