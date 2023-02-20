package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/services/v2ray"
	v2raytypes "github.com/sentinel-official/cli-client/services/v2ray/types"
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

			if err = status.LoadFromPath(statusFilePath); err != nil {
				return err
			}

			var service clienttypes.Service
			if status.Type == 1 {
				var cfg wireguardtypes.Config
				if err = json.Unmarshal(status.Info, &cfg); err != nil {
					return err
				}

				service = wireguard.NewWireGuard(&cfg)
			} else if status.Type == 2 {
				var cfg v2raytypes.Config
				if err = json.Unmarshal(status.Info, &cfg); err != nil {
					return err
				}

				service = v2ray.NewV2Ray(&cfg)
			} else {
				return nil
			}

			if service.IsUp() {
				if err = service.PreDown(); err != nil {
					return err
				}
				if err = service.Down(); err != nil {
					return err
				}
				if err = service.PostDown(); err != nil {
					return err
				}
			}

			return os.Remove(statusFilePath)
		},
	}

	return cmd
}
