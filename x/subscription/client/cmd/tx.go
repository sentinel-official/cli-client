package cmd

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/hub/x/subscription/types"
)

func GetTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscription",
		Short: "Subscription related subcommands",
	}

	cmd.AddCommand(
		txAllocate(),
		txCancel(),
	)

	return cmd
}

func txAllocate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allocate [id] [address] [bytes]",
		Short: "Add an allocation for subscription",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			address, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			bytes, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgAllocateRequest(
				ctx.FromAddress,
				id,
				address,
				sdk.NewInt(bytes),
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

func txCancel() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "cancel [id]",
		Short:  "Cancel a subscription",
		Args:   cobra.ExactArgs(1),
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelRequest(
				ctx.FromAddress,
				id,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
