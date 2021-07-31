package cmd

import (
	"context"
	"fmt"
	clienttypes "github.com/sentinel-official/cli-client/types"
	"github.com/sentinel-official/cli-client/utils"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	hubtypes "github.com/sentinel-official/hub/types"
	sessiontypes "github.com/sentinel-official/hub/x/session/types"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/x/session/types"
)

var (
	header = []string{
		"ID",
		"Subscription",
		"Node",
		"Address",
		"Duration",
		"Bandwidth",
		"Status",
	}
)

func QuerySession() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session [id]",
		Short: "Query a session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFormat, err := cmd.Flags().GetString(clienttypes.FlagOutput)
			if err != nil {
				return err
			}

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			var (
				qsc = sessiontypes.NewQueryServiceClient(ctx)
			)

			result, err := qsc.QuerySession(
				context.Background(),
				sessiontypes.NewQuerySessionRequest(id),
			)
			if err != nil {
				return err
			}

			var (
				item       = types.NewSessionFromRaw(&result.Session)
				outputRows [][]string
			)

			outputRows = append(
				outputRows,
				[]string{
					fmt.Sprintf("%d", item.ID),
					fmt.Sprintf("%d", item.Subscription),
					item.Node,
					item.Address,
					item.Duration.Truncate(1 * time.Second).String(),
					item.Bandwidth.String(),
					item.Status,
				},
			)

			err = utils.WriteOutput(header, outputRows, outputFormat)
			if err != nil {
				return err
			}
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QuerySessions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sessions",
		Short: "Query sessions",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			subscription, err := cmd.Flags().GetUint64(flagSubscription)
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

			outputFormat, err := cmd.Flags().GetString(clienttypes.FlagOutput)
			if err != nil {
				return err
			}

			var (
				items      types.Sessions
				qc         = sessiontypes.NewQueryServiceClient(ctx)
				outputRows [][]string
			)

			if subscription != 0 {
				result, err := qc.QuerySessionsForSubscription(
					context.Background(),
					sessiontypes.NewQuerySessionsForSubscriptionRequest(
						subscription,
						pagination,
					),
				)
				if err != nil {
					return err
				}

				items = append(items, types.NewSessionsFromRaw(result.Sessions)...)
			} else if bech32Address != "" {
				address, err := sdk.AccAddressFromBech32(bech32Address)
				if err != nil {
					return err
				}

				status, err := cmd.Flags().GetString(flagStatus)
				if err != nil {
					return err
				}

				result, err := qc.QuerySessionsForAddress(
					context.Background(),
					sessiontypes.NewQuerySessionsForAddressRequest(
						address,
						hubtypes.StatusFromString(status),
						pagination,
					),
				)
				if err != nil {
					return err
				}

				items = append(items, types.NewSessionsFromRaw(result.Sessions)...)
			} else {
				result, err := qc.QuerySessions(
					context.Background(),
					sessiontypes.NewQuerySessionsRequest(pagination),
				)
				if err != nil {
					return err
				}

				items = append(items, types.NewSessionsFromRaw(result.Sessions)...)
			}

			for i := 0; i < len(items); i++ {
				outputRows = append(
					outputRows,
					[]string{
						fmt.Sprintf("%d", items[i].ID),
						fmt.Sprintf("%d", items[i].Subscription),
						items[i].Node,
						items[i].Address,
						items[i].Duration.Truncate(1 * time.Second).String(),
						items[i].Bandwidth.String(),
						items[i].Status,
					},
				)
			}

			err = utils.WriteOutput(header, outputRows, outputFormat)
			if err != nil {
				return err
			}
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "sessions")

	cmd.Flags().String(flagAddress, "", "filter with account address")
	cmd.Flags().Uint64(flagSubscription, 0, "filter with subscription identity")
	cmd.Flags().String(flagStatus, "Active", "filter with status (Active|Inactive)")

	return cmd
}
