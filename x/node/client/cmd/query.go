package cmd

import (
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

	"github.com/sentinel-official/cli-client/context"
	clitypes "github.com/sentinel-official/cli-client/types"
	resttypes "github.com/sentinel-official/cli-client/types/rest"
	"github.com/sentinel-official/cli-client/x/node/types"
)

var (
	header = []string{
		"Moniker",
		"Address",
		"Provider",
		"Price",
		"Country",
		"Speed test",
		"Latency",
		"Peers",
		"Handshake",
		"Version",
		"Status",
	}
)

func fetchNodeInfo(remote string, timeout time.Duration) (info types.Info, err error) {
	var (
		body       resttypes.Response
		endpoint   = strings.Trim(remote, "/") + "/status"
		httpclient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: timeout,
		}
		startTime = time.Now()
	)

	resp, err := httpclient.Get(endpoint)
	if err != nil {
		return info, err
	}

	info.Latency = time.Since(startTime)
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
			qc, err := context.NewQueryContextFromCmd(cmd)
			if err != nil {
				return err
			}

			address, err := hubtypes.NodeAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			timeout, err := cmd.Flags().GetDuration(clitypes.FlagTimeout)
			if err != nil {
				return err
			}

			result, err := qc.QueryNode(address)
			if err != nil {
				return err
			}

			var (
				info, _ = fetchNodeInfo(result.RemoteURL, timeout)
				item    = types.NewNodeFromRaw(result).WithInfo(info)
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
					item.Bandwidth.String(),
					item.Latency.Truncate(1 * time.Millisecond).String(),
					fmt.Sprintf("%d", item.Peers),
					fmt.Sprintf("%t", item.Handshake.Enable),
					item.Version,
					item.Status,
				},
			)

			table.Render()
			return nil
		},
	}

	clitypes.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryNodes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Query nodes",
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

			timeout, err := cmd.Flags().GetDuration(clitypes.FlagTimeout)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			var (
				items  []nodetypes.Node
				status = hubtypes.StatusFromString(s)
			)

			if provider != "" {
				address, err := hubtypes.ProvAddressFromBech32(provider)
				if err != nil {
					return err
				}

				result, err := qc.QueryNodesForProvider(
					address,
					status,
					pagination,
				)
				if err != nil {
					return err
				}

				items = append(items, result...)
			} else {
				result, err := qc.QueryNodes(
					status,
					pagination,
				)
				if err != nil {
					return err
				}

				items = append(items, result...)
			}

			var (
				group = sync.WaitGroup{}
				mutex = sync.Mutex{}
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(header)
			for i := 0; i < len(items); i++ {
				group.Add(1)
				go func(i int) {
					defer group.Done()

					var (
						info, _ = fetchNodeInfo(items[i].RemoteURL, timeout)
						item    = types.NewNodeFromRaw(&items[i]).WithInfo(info)
					)

					mutex.Lock()
					defer mutex.Unlock()

					table.Append(
						[]string{
							item.Moniker,
							item.Address,
							item.Provider,
							item.Price.Raw().String(),
							item.Location.Country,
							item.Bandwidth.String(),
							item.Latency.Truncate(1 * time.Millisecond).String(),
							fmt.Sprintf("%d", item.Peers),
							fmt.Sprintf("%t", item.Handshake.Enable),
							item.Version,
							item.Status,
						},
					)
				}(i)
			}

			group.Wait()

			table.Render()
			return nil
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "nodes")
	clitypes.AddQueryFlagsToCmd(cmd)

	cmd.Flags().String(flagProvider, "", "filter with provider address")
	cmd.Flags().String(flagStatus, "Active", "filter with status (Active|Inactive)")

	return cmd
}
