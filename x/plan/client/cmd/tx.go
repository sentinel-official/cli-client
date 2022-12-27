package cmd

import (
	"bufio"
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	hubtypes "github.com/sentinel-official/hub/types"
	"github.com/sentinel-official/hub/x/plan/types"

	"github.com/sentinel-official/cli-client/context"
	clitypes "github.com/sentinel-official/cli-client/types"
)

func GetTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "plan",
		Short:                      "plan related subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		txAdd(),
		txSetStatus(),
		txAddNode(),
		txRemoveNode(),
	)

	return cmd
}

func txAdd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [price] [validity] [bytes]",
		Short: "Add a plan",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cc, err := context.NewClientContextFromCmd(cmd)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			password, from, err := cc.ReadPasswordAndGetAddress(reader, cc.From)
			if err != nil {
				return err
			}

			price, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			validity, err := time.ParseDuration(args[1])
			if err != nil {
				return err
			}

			bytes, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddRequest(
				from.Bytes(),
				price,
				validity,
				sdk.NewInt(bytes),
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

	return cmd
}

func txSetStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status-set [plan] [status]",
		Short: "Set a plan status",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cc, err := context.NewClientContextFromCmd(cmd)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			password, from, err := cc.ReadPasswordAndGetAddress(reader, cc.From)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetStatusRequest(
				from.Bytes(),
				id,
				hubtypes.StatusFromString(args[1]),
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

	return cmd
}

func txAddNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node-add [plan] [node]",
		Short: "Add a node for plan",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cc, err := context.NewClientContextFromCmd(cmd)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			password, from, err := cc.ReadPasswordAndGetAddress(reader, cc.From)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			node, err := hubtypes.NodeAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgAddNodeRequest(
				from.Bytes(),
				id,
				node,
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

	return cmd
}

func txRemoveNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node-remove [plan] [node]",
		Short: "Remove a node for plan",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cc, err := context.NewClientContextFromCmd(cmd)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			password, from, err := cc.ReadPasswordAndGetAddress(reader, cc.From)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			node, err := hubtypes.NodeAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveNodeRequest(
				from,
				id,
				node,
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

	return cmd
}
