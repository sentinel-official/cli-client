package cmd

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/olekukonko/tablewriter"
	deposittypes "github.com/sentinel-official/hub/x/deposit/types"
	"github.com/spf13/cobra"

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
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			var (
				qsc = deposittypes.NewQueryServiceClient(ctx)
			)

			result, err := qsc.QueryDeposit(
				context.Background(),
				deposittypes.NewQueryDepositRequest(address),
			)
			if err != nil {
				return err
			}

			var (
				item  = types.NewDepositFromRaw(&result.Deposit)
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

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryDeposits() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposits",
		Short: "Query deposits",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			var (
				qsc = deposittypes.NewQueryServiceClient(ctx)
			)

			result, err := qsc.QueryDeposits(
				context.Background(),
				deposittypes.NewQueryDepositsRequest(pagination),
			)
			if err != nil {
				return err
			}

			var (
				items = types.NewDepositsFromRaw(result.Deposits)
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

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "deposits")

	return cmd
}
