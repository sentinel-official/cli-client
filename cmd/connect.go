package cmd

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/pkg/errors"
	hubtypes "github.com/sentinel-official/hub/types"
	nodetypes "github.com/sentinel-official/hub/x/node/types"
	sessiontypes "github.com/sentinel-official/hub/x/session/types"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/services/wireguard"
	wireguardtypes "github.com/sentinel-official/cli-client/services/wireguard/types"
	clienttypes "github.com/sentinel-official/cli-client/types"
	netutil "github.com/sentinel-official/cli-client/utils/net"
)

func queryNode(qsc nodetypes.QueryServiceClient, address hubtypes.NodeAddress) (*nodetypes.Node, error) {
	var (
		result, err = qsc.QueryNode(
			context.Background(),
			nodetypes.NewQueryNodeRequest(address),
		)
	)

	if err != nil {
		return nil, err
	}

	return &result.Node, nil
}

func queryActiveSession(qsc sessiontypes.QueryServiceClient, address sdk.AccAddress) (*sessiontypes.Session, error) {
	var (
		result, err = qsc.QuerySessionsForAddress(
			context.Background(),
			sessiontypes.NewQuerySessionsForAddressRequest(
				address,
				hubtypes.StatusActive,
				&query.PageRequest{
					Limit: 1,
				},
			),
		)
	)

	if err != nil {
		return nil, err
	}
	if len(result.Sessions) > 0 {
		return &result.Sessions[0], nil
	}

	return nil, nil
}

func ConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect [subscription] [address]",
		Short: "Connect to a node",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			address, err := hubtypes.NodeAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			timeout, err := cmd.Flags().GetDuration(clienttypes.FlagTimeout)
			if err != nil {
				return err
			}

			ss, err := cmd.Flags().GetStringArray(clienttypes.FlagResolver)
			if err != nil {
				return err
			}

			var (
				resolvers      []net.IP
				status         = clienttypes.NewStatus()
				statusFilePath = filepath.Join(ctx.HomeDir, "status.json")
			)

			for _, s := range ss {
				ip := net.ParseIP(s)
				if ip == nil {
					return fmt.Errorf("provided resolver ip %s is invalid", s)
				}

				resolvers = append(resolvers, ip)
			}

			if err := status.LoadFromPath(statusFilePath); err != nil {
				return err
			}

			if status.IFace != "" {
				var (
					service = wireguard.NewWireGuard().
						WithConfig(
							&wireguardtypes.Config{
								Name: status.IFace,
							},
						)
				)

				if service.IsUp() {
					if err := service.PreDown(); err != nil {
						return err
					}
					if err := service.Down(); err != nil {
						return err
					}
					if err := service.PostDown(); err != nil {
						return err
					}
				}
			}

			var (
				messages           []sdk.Msg
				nodeQueryClient    = nodetypes.NewQueryServiceClient(ctx)
				sessionQueryClient = sessiontypes.NewQueryServiceClient(ctx)
			)

			session, err := queryActiveSession(sessionQueryClient, ctx.FromAddress)
			if err != nil {
				return err
			}

			if session != nil {
				messages = append(
					messages,
					sessiontypes.NewMsgEndRequest(
						ctx.FromAddress,
						session.Id,
						0,
					),
				)
			}

			messages = append(
				messages,
				sessiontypes.NewMsgStartRequest(
					ctx.FromAddress,
					id,
					address,
				),
			)

			if err := tx.GenerateOrBroadcastTxCLI(ctx, cmd.Flags(), messages...); err != nil {
				return err
			}

			node, err := queryNode(nodeQueryClient, address)
			if err != nil {
				return err
			}

			session, err = queryActiveSession(sessionQueryClient, ctx.FromAddress)
			if err != nil {
				return err
			}
			if session == nil {
				return errors.New("no active session found")
			}

			wgPrivateKey, err := wireguardtypes.NewPrivateKey()
			if err != nil {
				return err
			}

			signature, _, err := ctx.Keyring.Sign(ctx.From, sdk.Uint64ToBigEndian(session.Id))
			if err != nil {
				return err
			}

			req, err := json.Marshal(
				map[string]interface{}{
					"key":       wgPrivateKey.Public().String(),
					"signature": signature,
				},
			)
			if err != nil {
				return err
			}

			var (
				body     clienttypes.Response
				endpoint = fmt.Sprintf(
					"%s/accounts/%s/sessions/%d",
					strings.Trim(node.RemoteURL, "/"), ctx.FromAddress, session.Id,
				)
				httpClient = http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
					Timeout: timeout,
				}
			)

			resp, err := httpClient.Post(endpoint, jsonrpc.ContentType, bytes.NewBuffer(req))
			if err != nil {
				return err
			}

			defer resp.Body.Close()

			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				return err
			}
			if body.Error != nil {
				return errors.New(body.Error.Message)
			}

			result, err := base64.StdEncoding.DecodeString(body.Result.(string))
			if err != nil {
				return err
			}
			if len(result) != 58 {
				return fmt.Errorf("incorrect result size %d", len(result))
			}

			var (
				ipv4Address         = net.IP(result[0:4])
				ipv6Address         = net.IP(result[4:20])
				endpointHost        = net.IP(result[20:24])
				endpointPort        = binary.BigEndian.Uint16(result[24:26])
				endpointWGPublicKey = wireguardtypes.NewKey(result[26:58])
			)

			listenPort, err := netutil.GetFreeUDPPort()
			if err != nil {
				return err
			}

			var (
				cfg = &wireguardtypes.Config{
					Name: wireguardtypes.DefaultInterface,
					Interface: wireguardtypes.Interface{
						Addresses: []wireguardtypes.IPNet{
							{IP: ipv4Address, Net: 32},
							{IP: ipv6Address, Net: 128},
						},
						ListenPort: listenPort,
						PrivateKey: *wgPrivateKey,
						DNS: append(
							[]net.IP{net.ParseIP("10.8.0.1")},
							resolvers...,
						),
					},
					Peers: []wireguardtypes.Peer{
						{
							PublicKey: *endpointWGPublicKey,
							AllowedIPs: []wireguardtypes.IPNet{
								{IP: net.ParseIP("0.0.0.0")},
								{IP: net.ParseIP("::")},
							},
							Endpoint: wireguardtypes.Endpoint{
								Host: endpointHost.String(),
								Port: endpointPort,
							},
							PersistentKeepalive: 15,
						},
					},
				}

				service = wireguard.NewWireGuard().
					WithConfig(cfg)
			)

			status = clienttypes.NewStatus().
				WithFrom(ctx.FromAddress.String()).
				WithID(id).
				WithIFace(cfg.Name).
				WithTo(address.String())
			if err := status.SaveToPath(statusFilePath); err != nil {
				return err
			}

			if err := service.PreUp(); err != nil {
				return err
			}
			if err := service.Up(); err != nil {
				return err
			}
			if err := service.PostUp(); err != nil {
				return err
			}

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String(flags.FlagChainID, "", "the network chain identity")
	cmd.Flags().StringArray(clienttypes.FlagResolver, nil, "provide additional DNS servers")
	cmd.Flags().Duration(clienttypes.FlagTimeout, 15*time.Second, "time limit for requests made by the HTTP client")

	return cmd
}
