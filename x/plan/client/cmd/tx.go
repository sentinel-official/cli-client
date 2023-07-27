package cmd

import (
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	hubtypes "github.com/sentinel-official/hub/types"
	plantypes "github.com/sentinel-official/hub/x/plan/types"
)

func GetTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "plan related subcommands",
	}

	cmd.AddCommand(
		txCreate(),
		txUpdateStatus(),
		txLinkNode(),
		txUnlinkNode(),
		txSubscribe(),
	)

	return cmd
}

func txCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [duration] [gigabytes] [prices]",
		Short: "Create a plan",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			duration, err := time.ParseDuration(args[0])
			if err != nil {
				return err
			}

			gigabytes, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			prices, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}

			msg := plantypes.NewMsgCreateRequest(
				ctx.FromAddress.Bytes(),
				duration,
				gigabytes,
				prices,
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

func txUpdateStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-status [plan-id] [status]",
		Short: "Update a plan status",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			planID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := plantypes.NewMsgUpdateStatusRequest(
				ctx.FromAddress.Bytes(),
				planID,
				hubtypes.StatusFromString(args[1]),
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

func txLinkNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link-node [plan-id] [node-addr]",
		Short: "Link a node for plan",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			planID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			nodeAddr, err := hubtypes.NodeAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := plantypes.NewMsgLinkNodeRequest(
				ctx.FromAddress.Bytes(),
				planID,
				nodeAddr,
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

func txUnlinkNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unlink-node [plan-id] [node-addr]",
		Short: "Unlink a node for plan",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			planID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			nodeAddr, err := hubtypes.NodeAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := plantypes.NewMsgUnlinkNodeRequest(
				ctx.FromAddress.Bytes(),
				planID,
				nodeAddr,
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
		Use:   "subscribe [plan-id] [denom]",
		Short: "Subscribe to a plan",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := plantypes.NewMsgSubscribeRequest(
				ctx.FromAddress,
				id,
				args[1],
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
