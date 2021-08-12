package cmd

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/context"
	clitypes "github.com/sentinel-official/cli-client/types"
	"github.com/sentinel-official/cli-client/x/deposit/types"
)

var (
	header = []string{
		"Address",
		"Amount",
	}
)

func QueryDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [address]",
		Short: "Query a deposit",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			qc, err := context.NewQueryContextFromCmd(cmd)
			if err != nil {
				return err
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			result, err := qc.QueryDeposit(address)
			if err != nil {
				return err
			}

			var (
				item  = types.NewDepositFromRaw(result)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(header)
			table.Append(
				[]string{
					item.Address,
					item.Amount.Raw().String(),
				},
			)

			table.Render()
			return nil
		},
	}

	clitypes.AddQueryFlagsToCmd(cmd)
	_ = cmd.Flags().MarkHidden(clitypes.FlagTimeout)

	return cmd
}

func QueryDeposits() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposits",
		Short: "Query deposits",
		RunE: func(cmd *cobra.Command, args []string) error {
			qc, err := context.NewQueryContextFromCmd(cmd)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			result, err := qc.QueryDeposits(pagination)
			if err != nil {
				return err
			}

			var (
				items = types.NewDepositsFromRaw(result)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(header)
			for i := 0; i < len(items); i++ {
				table.Append(
					[]string{
						items[i].Address,
						items[i].Amount.Raw().String(),
					},
				)
			}

			table.Render()
			return nil
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "deposits")

	clitypes.AddQueryFlagsToCmd(cmd)
	_ = cmd.Flags().MarkHidden(clitypes.FlagTimeout)

	return cmd
}
