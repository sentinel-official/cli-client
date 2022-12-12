package main

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/sentinel-official/hub"
	hubtypes "github.com/sentinel-official/hub/types"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/cmd"
)

func main() {
	hubtypes.GetConfig().Seal()

	root := &cobra.Command{
		Use:                        "sentinelcli",
		SilenceUsage:               true,
		SuggestionsMinimumDistance: 2,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var (
				config = hub.MakeEncodingConfig()
				ctx    = client.Context{}.
					WithCodec(config.Marshaler).
					WithInterfaceRegistry(config.InterfaceRegistry).
					WithTxConfig(config.TxConfig).
					WithLegacyAmino(config.Amino)
			)

			if err := client.SetCmdClientContextHandler(ctx, cmd); err != nil {
				return err
			}

			return nil
		},
	}

	root.AddCommand(
		cmd.ConnectCmd(),
		cmd.DisconnectCmd(),
		cmd.KeysCmd(),
		cmd.QueryCommand(),
		cmd.StartCmd(),
		cmd.TxCommand(),
		version.NewVersionCommand(),
	)

	_ = root.ExecuteContext(
		context.WithValue(
			context.Background(),
			client.ClientContextKey,
			&client.Context{},
		),
	)
}
