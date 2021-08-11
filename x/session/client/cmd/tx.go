package cmd

import (
	"bufio"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	hubtypes "github.com/sentinel-official/hub/types"
	"github.com/sentinel-official/hub/x/session/types"

	"github.com/sentinel-official/cli-client/context"
	clitypes "github.com/sentinel-official/cli-client/types"
)

func GetTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "session",
		Short:                      "Session related subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		txStart(),
		txEnd(),
	)

	return cmd
}

func txStart() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [subscription] [node]",
		Short: "Start a session",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cc, err := context.NewClientContextFromCmd(cmd)
			if err != nil {
				return err
			}

			var (
				reader = bufio.NewReader(cmd.InOrStdin())
			)

			password, from, err := cc.ReadPasswordAndGetAddress(reader, cc.From)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			address, err := hubtypes.NodeAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgStartRequest(
				from,
				id,
				address,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			result, err := cc.SignAndBroadcastTx(password, msg)
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	}

	clitypes.AddTxFlagsToCmd(cmd)
	_ = cmd.Flags().MarkHidden(clitypes.FlagServiceHome)

	return cmd
}

func txEnd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "end [id]",
		Short: "End a session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cc, err := context.NewClientContextFromCmd(cmd)
			if err != nil {
				return err
			}

			var (
				reader = bufio.NewReader(cmd.InOrStdin())
			)

			password, from, err := cc.ReadPasswordAndGetAddress(reader, cc.From)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			rating, err := cmd.Flags().GetUint64(flagRating)
			if err != nil {
				return err
			}

			msg := types.NewMsgEndRequest(
				from,
				id,
				rating,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			result, err := cc.SignAndBroadcastTx(password, msg)
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	}

	clitypes.AddTxFlagsToCmd(cmd)
	_ = cmd.Flags().MarkHidden(clitypes.FlagServiceHome)

	cmd.Flags().Uint64(flagRating, 0, "rate the session quality [0, 10]")

	return cmd
}
