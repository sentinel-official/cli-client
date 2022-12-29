package cmd

import (
	"bufio"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/context"
	clitypes "github.com/sentinel-official/cli-client/types"
	cliutils "github.com/sentinel-official/cli-client/utils"
)

func KeysCmd() *cobra.Command {
	var (
		cmd = &cobra.Command{
			Use:                        "keys",
			Short:                      "Keys subcommands",
			DisableFlagParsing:         true,
			SuggestionsMinimumDistance: 2,
			RunE:                       client.ValidateCmd,
		}
	)

	cmd.AddCommand(
		addCmd(),
		deleteCmd(),
		listCmd(),
		showCmd(),
	)

	return cmd
}

func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "add [name]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			kc, err := context.NewKeyringContextFromCmd(cmd)
			if err != nil {
				return err
			}

			coinType, err := cmd.Flags().GetUint32(clitypes.FlagCoinType)
			if err != nil {
				return err
			}

			account, err := cmd.Flags().GetUint32(clitypes.FlagAccount)
			if err != nil {
				return err
			}

			index, err := cmd.Flags().GetUint32(clitypes.FlagIndex)
			if err != nil {
				return err
			}

			recoverKey, err := cmd.Flags().GetBool(clitypes.FlagRecover)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			password, err := cliutils.GetPassword(kc.Backend, reader)
			if err != nil {
				return err
			}

			entropy, err := bip39.NewEntropy(256)
			if err != nil {
				return err
			}

			mnemonic, err := bip39.NewMnemonic(entropy)
			if err != nil {
				return err
			}

			if recoverKey {
				mnemonic, err = input.GetString("Enter your bip39 mnemonic", reader)
				if err != nil {
					return err
				}
			}

			result, err := kc.AddKey(password, args[0], mnemonic, "", coinType, account, index)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "**Important** write this mnemonic phrase in a safe place\n")
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s\n", mnemonic)

			fmt.Println(result)
			return nil
		},
	}

	clitypes.AddKeyringFlagsToCmd(cmd)
	cmd.Flags().Bool(clitypes.FlagRecover, false, "provide mnemonic phrase to recover an existing key")
	cmd.Flags().Uint32(clitypes.FlagCoinType, 118, "coin type number for HD derivation")
	cmd.Flags().Uint32(clitypes.FlagAccount, 0, "account number for HD derivation")
	cmd.Flags().Uint32(clitypes.FlagIndex, 0, "address index number for HD derivation")

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete [name]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			kc, err := context.NewKeyringContextFromCmd(cmd)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			password, err := cliutils.GetPassword(kc.Backend, reader)
			if err != nil {
				return err
			}

			return kc.DeleteKey(password, args[0])
		},
	}

	clitypes.AddKeyringFlagsToCmd(cmd)

	return cmd
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			kc, err := context.NewKeyringContextFromCmd(cmd)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			password, err := cliutils.GetPassword(kc.Backend, reader)
			if err != nil {
				return err
			}

			result, err := kc.GetKeys(password)
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	}

	clitypes.AddKeyringFlagsToCmd(cmd)

	return cmd
}

func showCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "show [name]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			kc, err := context.NewKeyringContextFromCmd(cmd)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			password, err := cliutils.GetPassword(kc.Backend, reader)
			if err != nil {
				return err
			}

			result, err := kc.GetKey(password, args[0])
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	}

	clitypes.AddKeyringFlagsToCmd(cmd)

	return cmd
}
