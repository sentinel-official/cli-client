package cmd

import (
	"bufio"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	hubtypes "github.com/sentinel-official/hub/types"
	nodetypes "github.com/sentinel-official/hub/x/node/types"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/context"
	clitypes "github.com/sentinel-official/cli-client/types"
)

func GetTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "node",
		Short:                      "Node related subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
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
			tc, err := context.NewTxContextFromCmd(cmd)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			password, fromAddr, err := tc.GetPasswordAndAddress(reader, tc.From)
			if err != nil {
				return err
			}

			msg := nodetypes.NewMsgSetStatusRequest(
				fromAddr.Bytes(),
				hubtypes.StatusFromString(args[0]),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			result, err := tc.SignMessagesAndBroadcastTx(password, msg)
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
