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
			sc, err := context.NewServiceContextFromCmd(cmd)
			if err != nil {
				return err
			}

			status, err := sc.GetStatus()
			if err != nil {
				return err
			}

			if status.IFace != "" {
				if err := sc.Disconnect(); err != nil {
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
