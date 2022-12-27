package cmd

import (
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/context"
	clitypes "github.com/sentinel-official/cli-client/types"
)

func DisconnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect from a node",
		RunE: func(cmd *cobra.Command, args []string) error {
			cc, err := context.NewClientContextFromCmd(cmd)
			if err != nil {
				return err
			}

			status, err := cc.GetStatus()
			if err != nil {
				return err
			}

			if status.IFace != "" {
				if err := cc.Disconnect(); err != nil {
					return err
				}
			}

			return nil
		},
	}

	clitypes.AddServiceFlagsToCmd(cmd)
	clitypes.AddTimeoutFlagsToCmd(cmd)

	return cmd
}
