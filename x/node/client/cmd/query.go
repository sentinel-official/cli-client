package cmd

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/olekukonko/tablewriter"
	hubtypes "github.com/sentinel-official/hub/types"
	nodetypes "github.com/sentinel-official/hub/x/node/types"
	"github.com/spf13/cobra"

	clienttypes "github.com/sentinel-official/cli-client/types"
	"github.com/sentinel-official/cli-client/x/node/types"
)

var (
	header = []string{
		"Moniker",
		"Address",
		"Provider",
		"Price",
		"Country",
		"Peers",
		"Handshake",
		"Status",
	}
)

func fetchNodeInfo(remote string) (info types.Info, err error) {
	var (
		body       clienttypes.Response
		endpoint   = strings.Trim(remote, "/") + "/status"
		httpclient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: 5 * time.Second,
		}
	)

	resp, err := httpclient.Get(endpoint)
	if err != nil {
		return info, err
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return info, err
	}

	bytes, err := json.Marshal(body.Result)
	if err != nil {
		return info, err
	}

	if err := json.Unmarshal(bytes, &info); err != nil {
		return info, err
	}

	return info, nil
}

func QueryNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node [address]",
		Short: "Query a node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			address, err := hubtypes.NodeAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			var (
				qsc = nodetypes.NewQueryServiceClient(ctx)
			)

			result, err := qsc.QueryNode(
				context.Background(),
				nodetypes.NewQueryNodeRequest(address),
			)
			if err != nil {
				return err
			}

			var (
				info, _ = fetchNodeInfo(result.Node.RemoteURL)
				item    = types.NewNodeFromRaw(&result.Node).WithInfo(info)
				table   = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(header)
			table.Append(
				[]string{
					item.Moniker,
					item.Address,
					item.Provider,
					item.Price.Raw().String(),
					item.Location.Country,
					fmt.Sprintf("%d", item.Peers),
					fmt.Sprintf("%t", item.Handshake.Enable),
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

func QueryNodes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Query nodes",
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
				items  []nodetypes.Node
				qsc    = nodetypes.NewQueryServiceClient(ctx)
				status = hubtypes.StatusFromString(s)
			)

			if provider != "" {
				address, err := hubtypes.ProvAddressFromBech32(provider)
				if err != nil {
					return err
				}

				result, err := qsc.QueryNodesForProvider(
					context.Background(),
					nodetypes.NewQueryNodesForProviderRequest(address, status, pagination),
				)
				if err != nil {
					return err
				}

				items = append(items, result.Nodes...)
			} else {
				result, err := qsc.QueryNodes(
					context.Background(),
					nodetypes.NewQueryNodesRequest(status, pagination),
				)
				if err != nil {
					return err
				}

				items = append(items, result.Nodes...)
			}

			var (
				wg    = sync.WaitGroup{}
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(header)
			for i := 0; i < len(items); i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()

					var (
						info, _ = fetchNodeInfo(items[i].RemoteURL)
						item    = types.NewNodeFromRaw(&items[i]).WithInfo(info)
					)

					table.Append(
						[]string{
							item.Moniker,
							item.Address,
							item.Provider,
							item.Price.Raw().String(),
							item.Location.Country,
							fmt.Sprintf("%d", item.Peers),
							fmt.Sprintf("%t", item.Handshake.Enable),
							item.Status,
						},
					)
				}(i)
			}

			wg.Wait()

			table.Render()
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "nodes")

	cmd.Flags().String(flagProvider, "", "nodes operating under a provider")
	cmd.Flags().String(flagStatus, "Active", "nodes with status (Active|Inactive)")

	return cmd
}
