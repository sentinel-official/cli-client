package cmd

import (
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/services/wireguard"
	wireguardtypes "github.com/sentinel-official/cli-client/services/wireguard/types"
	clienttypes "github.com/sentinel-official/cli-client/types"
)

func DisconnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect from a node",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var (
				status         = clienttypes.NewStatus()
				statusFilePath = filepath.Join(ctx.HomeDir, "status.json")
			)

			if err := status.LoadFromPath(statusFilePath); err != nil {
				return err
			}

			if status.IFace != "" {
				var (
					service = wireguard.NewWireGuard().
						WithConfig(
							&wireguardtypes.Config{
								Name: status.IFace,
							},
						)
				)

				if service.IsUp() {
					if err := service.PreDown(); err != nil {
						return err
					}
					if err := service.Down(); err != nil {
						return err
					}
					if err := service.PostDown(); err != nil {
						return err
					}
				}

				return os.Remove(statusFilePath)
			}

			return nil
		},
	}

	return cmd
}
