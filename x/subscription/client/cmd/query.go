package cmd

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/olekukonko/tablewriter"
	hubtypes "github.com/sentinel-official/hub/types"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/context"
	clitypes "github.com/sentinel-official/cli-client/types"
	netutils "github.com/sentinel-official/cli-client/utils/net"
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
			qc, err := context.NewQueryContextFromCmd(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			result, err := qc.QuerySubscription(id)
			if err != nil {
				return err
			}

			var (
				item  = types.NewSubscriptionFromRaw(result)
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
					netutils.ToReadable(item.Free, 2),
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

func QuerySubscriptions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "Query subscriptions",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			qc, err := context.NewQueryContextFromCmd(cmd)
			if err != nil {
				return err
			}

			bech32Address, err := cmd.Flags().GetString(flagAddress)
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
			)

			if bech32Address != "" {
				address, err := sdk.AccAddressFromBech32(bech32Address)
				if err != nil {
					return err
				}

				result, err := qc.QuerySubscriptionsForAddress(
					address,
					hubtypes.StatusFromString(status),
					pagination,
				)
				if err != nil {
					return err
				}

				items = append(items, types.NewSubscriptionsFromRaw(result)...)
			} else {
				result, err := qc.QuerySubscriptions(pagination)
				if err != nil {
					return err
				}

				items = append(items, types.NewSubscriptionsFromRaw(result)...)
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
						netutils.ToReadable(items[i].Free, 2),
						items[i].Status,
					},
				)
			}

			table.Render()
			return nil
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "subscriptions")

	clitypes.AddQueryFlagsToCmd(cmd)
	_ = cmd.Flags().MarkHidden(clitypes.FlagTimeout)

	cmd.Flags().String(flagAddress, "", "filter with account address")
	cmd.Flags().String(flagStatus, "Active", "filter with status (Active|Inactive)")

	return cmd
}

func QueryQuota() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quota [subscription] [address]",
		Short: "Query a quota",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			qc, err := context.NewQueryContextFromCmd(cmd)
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

			result, err := qc.QueryQuota(
				id,
				address,
			)
			if err != nil {
				return err
			}

			var (
				item  = types.NewQuotaFromRaw(result)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(quotaHeader)
			table.Append(
				[]string{
					item.Address,
					netutils.ToReadable(item.Allocated, 2),
					netutils.ToReadable(item.Consumed, 2),
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

func QueryQuotas() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quotas [subscription]",
		Short: "Query quotas of a subscription",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			qc, err := context.NewQueryContextFromCmd(cmd)
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

			result, err := qc.QueryQuotas(
				id,
				pagination,
			)
			if err != nil {
				return err
			}

			var (
				items = types.NewQuotasFromRaw(result)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(quotaHeader)
			for i := 0; i < len(items); i++ {
				table.Append(
					[]string{
						items[i].Address,
						netutils.ToReadable(items[i].Allocated, 2),
						netutils.ToReadable(items[i].Consumed, 2),
					},
				)
			}

			table.Render()
			return nil
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "quotas")

	clitypes.AddQueryFlagsToCmd(cmd)
	_ = cmd.Flags().MarkHidden(clitypes.FlagTimeout)

	return cmd
}
