package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/version"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sentinel-official/hub"
	hubtypes "github.com/sentinel-official/hub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sentinel-official/cli-client/cmd"
	"github.com/sentinel-official/cli-client/types"
	configtypes "github.com/sentinel-official/cli-client/types/config"
)

func main() {
	hubtypes.GetConfig().Seal()
	root := &cobra.Command{
		Use:          "sentinelcli",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var (
				home       = viper.GetString(flags.FlagHome)
				configPath = filepath.Join(home, "config.toml")
			)

			if _, err := os.Stat(configPath); err != nil {
				if err := os.MkdirAll(home, 0700); err != nil {
					return err
				}

				config := configtypes.NewConfig().
					WithDefaultValues()
				if err := config.SaveToPath(configPath); err != nil {
					return err
				}
			}

			var (
				config    = hub.MakeEncodingConfig()
				clientCtx = client.Context{}.
					WithJSONMarshaler(config.Marshaler).
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
		cmd.StartCmd(),
		cmd.TxCommand(),
		keys.Commands(types.DefaultHomeDirectory),
		version.NewVersionCommand(),
	)

	root.PersistentFlags().String(flags.FlagHome, types.DefaultHomeDirectory, "home directory of the application")
	_ = viper.BindPFlag(flags.FlagHome, root.PersistentFlags().Lookup(flags.FlagHome))

	_ = root.ExecuteContext(
		context.WithValue(
			context.Background(),
			client.ClientContextKey,
			&client.Context{},
		),
	)
}
