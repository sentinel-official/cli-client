package cmd

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/hub/x/provider/types"
)

func GetTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider",
		Short: "Provider related subcommands",
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
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			identity, err := cmd.Flags().GetString(flagIdentity)
			if err != nil {
				return err
			}

			website, err := cmd.Flags().GetString(flagWebsite)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(flagDescription)
			if err != nil {
				return err
			}

			msg := types.NewMsgRegisterRequest(
				ctx.FromAddress,
				args[0],
				identity,
				website,
				description,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(flagIdentity, "", "identity signature (optional)")
	cmd.Flags().String(flagWebsite, "", "website (optional)")
	cmd.Flags().String(flagDescription, "", "description (optional)")

	return cmd
}

func txUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString(flagName)
			if err != nil {
				return err
			}

			identity, err := cmd.Flags().GetString(flagIdentity)
			if err != nil {
				return err
			}

			website, err := cmd.Flags().GetString(flagWebsite)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(flagDescription)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateRequest(
				ctx.FromAddress.Bytes(),
				name,
				identity,
				website,
				description,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(flagName, "", "name (optional)")
	cmd.Flags().String(flagIdentity, "", "identity signature (optional)")
	cmd.Flags().String(flagWebsite, "", "website (optional)")
	cmd.Flags().String(flagDescription, "", "description (optional)")

	return cmd
}
