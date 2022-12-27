package cmd

import (
	"bufio"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/hub/x/provider/types"

	"github.com/sentinel-official/cli-client/context"
	clitypes "github.com/sentinel-official/cli-client/types"
)

func GetTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "provider",
		Short:                      "Provider related subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		txRegister(),
		txUpdate(),
	)

	return cmd
}

func txRegister() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "register [name]",
		Short:  "Register a provider",
		Args:   cobra.ExactArgs(1),
		Hidden: true,
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

			identity, err := cmd.Flags().GetString(clitypes.FlagIdentity)
			if err != nil {
				return err
			}

			website, err := cmd.Flags().GetString(clitypes.FlagWebsite)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(clitypes.FlagDescription)
			if err != nil {
				return err
			}

			msg := types.NewMsgRegisterRequest(
				from,
				args[0],
				identity,
				website,
				description,
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

	cmd.Flags().String(clitypes.FlagIdentity, "", "identity signature (optional)")
	cmd.Flags().String(clitypes.FlagWebsite, "", "website (optional)")
	cmd.Flags().String(clitypes.FlagDescription, "", "description (optional)")

	return cmd
}

func txUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a provider",
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

			name, err := cmd.Flags().GetString(clitypes.FlagName)
			if err != nil {
				return err
			}

			identity, err := cmd.Flags().GetString(clitypes.FlagIdentity)
			if err != nil {
				return err
			}

			website, err := cmd.Flags().GetString(clitypes.FlagWebsite)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(clitypes.FlagDescription)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateRequest(
				from.Bytes(),
				name,
				identity,
				website,
				description,
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

	cmd.Flags().String(clitypes.FlagName, "", "name (optional)")
	cmd.Flags().String(clitypes.FlagIdentity, "", "identity signature (optional)")
	cmd.Flags().String(clitypes.FlagWebsite, "", "website (optional)")
	cmd.Flags().String(clitypes.FlagDescription, "", "description (optional)")

	return cmd
}
