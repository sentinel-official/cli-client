package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
	hubtypes "github.com/sentinel-official/hub/types"
	nodetypes "github.com/sentinel-official/hub/x/node/types"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/context"
	clitypes "github.com/sentinel-official/cli-client/types"
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

func fetchNodeInfo(remote string, timeout time.Duration) (info types.NodeInfo, err error) {
	var (
		body   clitypes.RestResponseBody
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: timeout,
		}
	)

	path, err := url.JoinPath(remote, "status")
	if err != nil {
		return info, err
	}

	start := time.Now()

	resp, err := client.Get(path)
	if err != nil {
		return info, err
	}

	info.Latency = time.Since(start)
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return info, err
	}

	result, err := json.Marshal(body.Result)
	if err != nil {
		return info, err
	}

	if err := json.Unmarshal(result, &info); err != nil {
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

			nodeAddr, err := hubtypes.NodeAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			timeout, err := cmd.Flags().GetDuration(clitypes.FlagTimeout)
			if err != nil {
				return err
			}

			result, err := qc.QueryNode(nodeAddr)
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

			provAddr, err := clitypes.GetProvAddressFromCmd(cmd)
			if err != nil {
				return err
			}

			status, err := clitypes.GetStatusFromCmd(cmd)
			if err != nil {
				return err
			}

			timeout, err := cmd.Flags().GetDuration(clitypes.FlagTimeout)
			if err != nil {
				return err
			}

			pagination, err := clitypes.GetPageRequestFromCmd(cmd)
			if err != nil {
				return err
			}

			var items []nodetypes.Node
			if provAddr != nil {
				result, err := qc.QueryNodesForProvider(
					provAddr,
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
				wg    = sync.WaitGroup{}
				mutex = sync.Mutex{}
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(header)
			for i := 0; i < len(items); i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()

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

			wg.Wait()

			table.Render()
			return nil
		},
	}

	clitypes.AddQueryFlagsToCmd(cmd)
	clitypes.AddPaginationFlagsToCmd(cmd, "nodes")

	cmd.Flags().String(clitypes.FlagProvider, "", "filter with provider address")
	cmd.Flags().String(clitypes.FlagStatus, "Active", "filter with status (Active|Inactive)")

	return cmd
}
