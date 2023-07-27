package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/olekukonko/tablewriter"
	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"
	"github.com/spf13/cobra"

	netutil "github.com/sentinel-official/cli-client/utils/net"
	"github.com/sentinel-official/cli-client/x/subscription/types"
)

var (
	subscriptionHeader = []string{
		"ID",
		"Address",
		"Inactive at",
		"Status",
		"Node",
		"Gigabytes",
		"Hours",
		"Deposit",
		"Plan",
		"Denom",
	}
	allocationHeader = []string{
		"Address",
		"Granted bytes",
		"Utilised bytes",
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

			var subscription subscriptiontypes.Subscription
			if err = ctx.InterfaceRegistry.UnpackAny(result.Subscription, &subscription); err != nil {
				return err
			}

			var (
				item  = types.NewSubscriptionFromRaw(subscription)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(subscriptionHeader)
			table.Append(
				[]string{
					fmt.Sprintf("%d", item.ID),
					item.Address,
					item.InactiveAt.String(),
					item.Status,
					item.NodeAddress,
					fmt.Sprintf("%d", item.Gigabytes),
					fmt.Sprintf("%d", item.Hours),
					item.Deposit,
					fmt.Sprintf("%d", item.PlanID),
					item.Denom,
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

				result, err := qsc.QuerySubscriptionsForAccount(
					context.Background(),
					subscriptiontypes.NewQuerySubscriptionsForAccountRequest(
						address,
						pagination,
					),
				)
				if err != nil {
					return err
				}

				var subscriptions []subscriptiontypes.Subscription
				for _, item := range result.Subscriptions {
					var subscription subscriptiontypes.Subscription
					if err = ctx.InterfaceRegistry.UnpackAny(item, &subscription); err != nil {
						return err
					}

					subscriptions = append(subscriptions, subscription)
				}

				items = append(items, types.NewSubscriptionsFromRaw(subscriptions)...)
			} else {
				result, err := qsc.QuerySubscriptions(
					context.Background(),
					subscriptiontypes.NewQuerySubscriptionsRequest(pagination),
				)
				if err != nil {
					return err
				}

				var subscriptions []subscriptiontypes.Subscription
				for _, item := range result.Subscriptions {
					var subscription subscriptiontypes.Subscription
					if err = ctx.InterfaceRegistry.UnpackAny(item, &subscription); err != nil {
						return err
					}

					subscriptions = append(subscriptions, subscription)
				}

				items = append(items, types.NewSubscriptionsFromRaw(subscriptions)...)
			}

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader(subscriptionHeader)

			for i := 0; i < len(items); i++ {
				table.Append(
					[]string{
						fmt.Sprintf("%d", items[i].ID),
						items[i].Address,
						items[i].InactiveAt.String(),
						items[i].Status,
						items[i].NodeAddress,
						fmt.Sprintf("%d", items[i].Gigabytes),
						fmt.Sprintf("%d", items[i].Hours),
						items[i].Deposit,
						fmt.Sprintf("%d", items[i].PlanID),
						items[i].Denom,
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

	return cmd
}

func QueryAllocation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allocation [subscription] [address]",
		Short: "Query a allocation",
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

			result, err := qsc.QueryAllocation(
				context.Background(),
				subscriptiontypes.NewQueryAllocationRequest(
					id,
					address,
				),
			)
			if err != nil {
				return err
			}

			var (
				item  = types.NewAllocationFromRaw(&result.Allocation)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(allocationHeader)
			table.Append(
				[]string{
					item.Address,
					netutil.ToReadable(item.GrantedBytes, 2),
					netutil.ToReadable(item.UtilisedBytes, 2),
				},
			)

			table.Render()
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryAllocations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allocations [subscription]",
		Short: "Query allocations of a subscription",
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

			result, err := qsc.QueryAllocations(
				context.Background(),
				subscriptiontypes.NewQueryAllocationsRequest(
					id,
					pagination,
				),
			)
			if err != nil {
				return err
			}

			var (
				items = types.NewAllocationsFromRaw(result.Allocations)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(allocationHeader)
			for i := 0; i < len(items); i++ {
				table.Append(
					[]string{
						items[i].Address,
						netutil.ToReadable(items[i].GrantedBytes, 2),
						netutil.ToReadable(items[i].UtilisedBytes, 2),
					},
				)
			}

			table.Render()
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "allocations")

	return cmd
}
