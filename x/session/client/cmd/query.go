package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/olekukonko/tablewriter"
	hubtypes "github.com/sentinel-official/hub/types"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/context"
	clitypes "github.com/sentinel-official/cli-client/types"
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
			qc, err := context.NewQueryContextFromCmd(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			result, err := qc.QuerySession(id)
			if err != nil {
				return err
			}

			var (
				item  = types.NewSessionFromRaw(result)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(header)
			table.Append(
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

			table.Render()
			return nil
		},
	}

	clitypes.AddQueryFlagsToCmd(cmd)
	_ = cmd.Flags().MarkHidden(clitypes.FlagTimeout)

	return cmd
}

func QuerySessions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sessions",
		Short: "Query sessions",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			qc, err := context.NewQueryContextFromCmd(cmd)
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
				items types.Sessions
			)

			if bech32Address != "" {
				address, err := sdk.AccAddressFromBech32(bech32Address)
				if err != nil {
					return err
				}

				status, err := cmd.Flags().GetString(flagStatus)
				if err != nil {
					return err
				}

				result, err := qc.QuerySessionsForAddress(
					address,
					hubtypes.StatusFromString(status),
					pagination,
				)
				if err != nil {
					return err
				}

				items = append(items, types.NewSessionsFromRaw(result)...)
			} else {
				result, err := qc.QuerySessions(
					pagination,
				)
				if err != nil {
					return err
				}

				items = append(items, types.NewSessionsFromRaw(result)...)
			}

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader(header)

			for i := 0; i < len(items); i++ {
				table.Append(
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

			table.Render()
			return nil
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "sessions")

	clitypes.AddQueryFlagsToCmd(cmd)
	_ = cmd.Flags().MarkHidden(clitypes.FlagTimeout)

	cmd.Flags().String(flagAddress, "", "filter with account address")
	cmd.Flags().String(flagStatus, "Active", "filter with status (Active|Inactive)")

	return cmd
}
