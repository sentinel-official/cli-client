package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/olekukonko/tablewriter"
	hubtypes "github.com/sentinel-official/hub/types"
	plantypes "github.com/sentinel-official/hub/x/plan/types"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/x/plan/types"
)

var (
	header = []string{
		"ID",
		"Provider",
		"Price",
		"Bytes",
		"Validity",
		"Status",
	}
)

func QueryPlan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan [id]",
		Short: "Query a plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			var (
				qsc = plantypes.NewQueryServiceClient(ctx)
			)

			result, err := qsc.QueryPlan(
				context.Background(),
				plantypes.NewQueryPlanRequest(id),
			)
			if err != nil {
				return err
			}

			var (
				item  = types.NewPlanFromRaw(&result.Plan)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(header)
			table.Append(
				[]string{
					fmt.Sprintf("%d", item.ID),
					item.Provider,
					item.Price.Raw().String(),
					fmt.Sprintf("%d", item.Bytes),
					item.Validity.String(),
					item.Status,
				},
			)

			table.Render()
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryPlans() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plans",
		Short: "Query plans",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			provider, err := cmd.Flags().GetString(flagProvider)
			if err != nil {
				return err
			}

			s, err := cmd.Flags().GetString(flagStatus)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			var (
				items  types.Plans
				qsc    = plantypes.NewQueryServiceClient(ctx)
				status = hubtypes.StatusFromString(s)
			)

			if provider != "" {
				address, err := hubtypes.ProvAddressFromBech32(provider)
				if err != nil {
					return err
				}

				result, err := qsc.QueryPlansForProvider(
					context.Background(),
					plantypes.NewQueryPlansForProviderRequest(
						address,
						status,
						pagination,
					),
				)
				if err != nil {
					return err
				}

				items = append(items, types.NewPlansFromRaw(result.Plans)...)
			} else {
				result, err := qsc.QueryPlans(
					context.Background(),
					plantypes.NewQueryPlansRequest(
						status,
						pagination,
					),
				)
				if err != nil {
					return err
				}

				items = append(items, types.NewPlansFromRaw(result.Plans)...)
			}

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader(header)

			for i := 0; i < len(items); i++ {
				table.Append(
					[]string{
						fmt.Sprintf("%d", items[i].ID),
						items[i].Provider,
						items[i].Price.Raw().String(),
						fmt.Sprintf("%d", items[i].Bytes),
						items[i].Validity.String(),
						items[i].Status,
					},
				)
			}

			table.Render()
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "plans")

	cmd.Flags().String(flagProvider, "", "filter with provider address")
	cmd.Flags().String(flagStatus, "Active", "filter with status")

	return cmd
}
