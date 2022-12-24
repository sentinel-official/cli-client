package cmd

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/olekukonko/tablewriter"
	hubtypes "github.com/sentinel-official/hub/types"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/context"
	clitypes "github.com/sentinel-official/cli-client/types"
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
			qc, err := context.NewQueryContextFromCmd(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			result, err := qc.QueryPlan(id)
			if err != nil {
				return err
			}

			var (
				item  = types.NewPlanFromRaw(result)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(header)
			table.Append(
				[]string{
					fmt.Sprintf("%d", item.ID),
					item.Provider,
					item.Price.Raw().String(),
					clitypes.ToReadableBytes(item.Bytes, 2),
					item.Validity.String(),
					item.Status,
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

func QueryPlans() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plans",
		Short: "Query plans",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			qc, err := context.NewQueryContextFromCmd(cmd)
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
				status = hubtypes.StatusFromString(s)
			)

			if provider != "" {
				address, err := hubtypes.ProvAddressFromBech32(provider)
				if err != nil {
					return err
				}

				result, err := qc.QueryPlansForProvider(
					address,
					status,
					pagination,
				)
				if err != nil {
					return err
				}

				items = append(items, types.NewPlansFromRaw(result)...)
			} else {
				result, err := qc.QueryPlans(
					status,
					pagination,
				)
				if err != nil {
					return err
				}

				items = append(items, types.NewPlansFromRaw(result)...)
			}

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader(header)

			for i := 0; i < len(items); i++ {
				table.Append(
					[]string{
						fmt.Sprintf("%d", items[i].ID),
						items[i].Provider,
						items[i].Price.Raw().String(),
						clitypes.ToReadableBytes(items[i].Bytes, 2),
						items[i].Validity.String(),
						items[i].Status,
					},
				)
			}

			table.Render()
			return nil
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "plans")

	clitypes.AddQueryFlagsToCmd(cmd)
	_ = cmd.Flags().MarkHidden(clitypes.FlagTimeout)

	cmd.Flags().String(flagProvider, "", "filter with provider address")
	cmd.Flags().String(flagStatus, "Active", "filter with status")

	return cmd
}
