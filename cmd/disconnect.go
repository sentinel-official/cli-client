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

	clitypes.AddFlagsToCmd(cmd)
	_ = cmd.Flags().MarkHidden(clitypes.FlagBroadcastMode)
	_ = cmd.Flags().MarkHidden(clitypes.FlagChainID)
	_ = cmd.Flags().MarkHidden(clitypes.FlagFrom)
	_ = cmd.Flags().MarkHidden(clitypes.FlagGas)
	_ = cmd.Flags().MarkHidden(clitypes.FlagGasPrices)
	_ = cmd.Flags().MarkHidden(clitypes.FlagKeyringBackend)
	_ = cmd.Flags().MarkHidden(clitypes.FlagKeyringHome)
	_ = cmd.Flags().MarkHidden(clitypes.FlagMemo)
	_ = cmd.Flags().MarkHidden(clitypes.FlagRPCAddress)

	return cmd
}
