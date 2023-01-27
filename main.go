package main

import (
	"context"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/version"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	hubtypes "github.com/sentinel-official/hub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sentinel-official/cli-client/cmd"
	"github.com/sentinel-official/cli-client/types"
)

func main() {
	hubtypes.GetConfig().Seal()
	root := &cobra.Command{
		Use:          "sentinelcli",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			home := viper.GetString(flags.FlagHome)
			if _, err := os.Stat(home); err != nil {
				if err := os.MkdirAll(home, 0700); err != nil {
					return err
				}
			}

			var (
				config    = types.MakeEncodingConfig()
				clientCtx = client.Context{}.
						WithCodec(config.Marshaler).
						WithInterfaceRegistry(config.InterfaceRegistry).
						WithTxConfig(config.TxConfig).
						WithLegacyAmino(config.Amino).
						WithInput(os.Stdin).
						WithAccountRetriever(authtypes.AccountRetriever{}).
						WithBroadcastMode(flags.BroadcastBlock).
						WithHomeDir(types.DefaultHomeDirectory)
			)

			if err := client.SetCmdClientContextHandler(clientCtx, cmd); err != nil {
				return err
			}

			return nil
		},
	}

	root.AddCommand(
		cmd.ConnectCmd(),
		cmd.DisconnectCmd(),
		cmd.QueryCommand(),
		cmd.TxCommand(),
		keys.Commands(types.DefaultHomeDirectory),
		version.NewVersionCommand(),
	)

	root.PersistentFlags().String(flags.FlagHome, types.DefaultHomeDirectory, "application home directory")
	_ = viper.BindPFlag(flags.FlagHome, root.PersistentFlags().Lookup(flags.FlagHome))

	_ = root.ExecuteContext(
		context.WithValue(
			context.Background(),
			client.ClientContextKey,
			&client.Context{},
		),
	)
}
