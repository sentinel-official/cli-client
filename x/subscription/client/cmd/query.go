package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/olekukonko/tablewriter"
	hubtypes "github.com/sentinel-official/hub/types"
	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/x/subscription/types"
)

var (
	subscriptionHeader = []string{
		"ID",
		"Owner",
		"Plan",
		"Expiry",
		"Denom",
		"Node",
		"Price",
		"Deposit",
		"Free",
		"Status",
	}
	quotaHeader = []string{
		"Address",
		"Allocated",
		"Consumed",
	}
)

func QuerySubscription() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscription [id]",
		Short: "Query a subscription",
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
				qsc = subscriptiontypes.NewQueryServiceClient(ctx)
			)

			result, err := qsc.QuerySubscription(
				context.Background(),
				subscriptiontypes.NewQuerySubscriptionRequest(id),
			)
			if err != nil {
				return err
			}

			var (
				item  = types.NewSubscriptionFromRaw(&result.Subscription)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(subscriptionHeader)
			table.Append(
				[]string{
					fmt.Sprintf("%d", item.ID),
					item.Owner,
					fmt.Sprintf("%d", item.Plan),
					item.Expiry.String(),
					item.Denom,
					item.Node,
					item.Price.Raw().String(),
					item.Deposit.Raw().String(),
					fmt.Sprintf("%d", item.Free),
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

func QuerySubscriptions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "Query subscriptions",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			bech32Address, err := cmd.Flags().GetString(flagAddress)
			if err != nil {
				return err
			}

			plan, err := cmd.Flags().GetUint64(flagPlan)
			if err != nil {
				return err
			}

			status, err := cmd.Flags().GetString(flagStatus)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			var (
				items types.Subscriptions
				qsc   = subscriptiontypes.NewQueryServiceClient(ctx)
			)

			if bech32Address != "" {
				address, err := sdk.AccAddressFromBech32(bech32Address)
				if err != nil {
					return err
				}

				result, err := qsc.QuerySubscriptionsForAddress(
					context.Background(),
					subscriptiontypes.NewQuerySubscriptionsForAddressRequest(
						address,
						hubtypes.StatusFromString(status),
						pagination,
					),
				)
				if err != nil {
					return err
				}

				items = append(items, types.NewSubscriptionsFromRaw(result.Subscriptions)...)
			} else if plan != 0 {
				result, err := qsc.QuerySubscriptionsForPlan(
					context.Background(),
					subscriptiontypes.NewQuerySubscriptionsForPlanRequest(
						plan,
						pagination,
					),
				)
				if err != nil {
					return err
				}

				items = append(items, types.NewSubscriptionsFromRaw(result.Subscriptions)...)
			} else {
				result, err := qsc.QuerySubscriptions(
					context.Background(),
					subscriptiontypes.NewQuerySubscriptionsRequest(pagination),
				)
				if err != nil {
					return err
				}

				items = append(items, types.NewSubscriptionsFromRaw(result.Subscriptions)...)
			}

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader(subscriptionHeader)

			for i := 0; i < len(items); i++ {
				table.Append(
					[]string{
						fmt.Sprintf("%d", items[i].ID),
						items[i].Owner,
						fmt.Sprintf("%d", items[i].Plan),
						items[i].Expiry.String(),
						items[i].Denom,
						items[i].Node,
						items[i].Price.Raw().String(),
						items[i].Deposit.Raw().String(),
						fmt.Sprintf("%d", items[i].Free),
						items[i].Status,
					},
				)
			}

			table.Render()
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "subscriptions")

	cmd.Flags().String(flagAddress, "", "filter with account address")
	cmd.Flags().Uint64(flagPlan, 0, "filter with plan identity")
	cmd.Flags().String(flagStatus, "Active", "filter with status (Active|Inactive)")

	return cmd
}

func QueryQuota() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quota [subscription] [address]",
		Short: "Query a quota",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			address, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			var (
				qsc = subscriptiontypes.NewQueryServiceClient(ctx)
			)

			result, err := qsc.QueryQuota(
				context.Background(),
				subscriptiontypes.NewQueryQuotaRequest(
					id,
					address,
				),
			)
			if err != nil {
				return err
			}

			var (
				item  = types.NewQuotaFromRaw(&result.Quota)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(quotaHeader)
			table.Append(
				[]string{
					item.Address,
					fmt.Sprintf("%d", item.Allocated),
					fmt.Sprintf("%d", item.Consumed),
				},
			)

			table.Render()
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryQuotas() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quotas [subscription]",
		Short: "Query quotas of a subscription",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			var (
				qsc = subscriptiontypes.NewQueryServiceClient(ctx)
			)

			result, err := qsc.QueryQuotas(
				context.Background(),
				subscriptiontypes.NewQueryQuotasRequest(
					id,
					pagination,
				),
			)
			if err != nil {
				return err
			}

			var (
				items = types.NewQuotasFromRaw(result.Quotas)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(quotaHeader)
			for i := 0; i < len(items); i++ {
				table.Append(
					[]string{
						items[i].Address,
						fmt.Sprintf("%d", items[i].Allocated),
						fmt.Sprintf("%d", items[i].Consumed),
					},
				)
			}

			table.Render()
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "quotas")

	return cmd
}
