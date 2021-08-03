package service

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/sentinel-official/cli-client/context"
	"github.com/sentinel-official/cli-client/services/wireguard"
	wireguardtypes "github.com/sentinel-official/cli-client/services/wireguard/types"
	"github.com/sentinel-official/cli-client/types"
	resttypes "github.com/sentinel-official/cli-client/types/rest"
	netutils "github.com/sentinel-official/cli-client/utils/net"
	restutils "github.com/sentinel-official/cli-client/utils/rest"
)

func HandlerConnect(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			config         = ctx.Config()
			status         = types.NewStatus()
			statusFilePath = filepath.Join(ctx.Home(), types.StatusFilename)
		)

		request, err := NewRequestConnect(r)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusBadRequest,
				resttypes.NewError(1001, err.Error()),
			)
			return
		}
		if err := request.Validate(); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusBadRequest,
				resttypes.NewError(1002, err.Error()),
			)
			return
		}

		if err := status.LoadFromPath(statusFilePath); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1003, err.Error()),
			)
			return
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
				restutils.WriteErrorToResponse(
					w, http.StatusBadRequest,
					resttypes.NewError(1004, fmt.Sprintf("service is already running in interface %s", status.IFace)),
				)
				return
			}
		}

		kr, err := keyring.New(
			version.Name,
			config.Keyring.Backend,
			ctx.Home(),
			strings.NewReader(""),
		)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1005, err.Error()),
			)
			return
		}

		key, err := kr.Key(request.From)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1006, err.Error()),
			)
			return
		}

		listenPort, err := netutils.GetFreeUDPPort()
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1007, err.Error()),
			)
			return
		}

		var (
			wireGuardConfig = &wireguardtypes.Config{
				Name: wireguardtypes.DefaultInterface,
				Interface: wireguardtypes.Interface{
					Addresses: []wireguardtypes.IPNet{
						{
							IP:  net.IP(request.Info[0 : 0+4]),
							Net: 32,
						},
						{
							IP:  net.IP(request.Info[4 : 4+16]),
							Net: 128,
						},
					},
					ListenPort: listenPort,
					PrivateKey: *wireguardtypes.NewKey(request.Keys[0]),
					DNS: append(
						[]net.IP{net.ParseIP("10.8.0.1")},
						request.Resolvers...,
					),
				},
				Peers: []wireguardtypes.Peer{
					{
						PublicKey: *wireguardtypes.NewKey(request.Info[26 : 26+32]),
						AllowedIPs: []wireguardtypes.IPNet{
							{IP: net.ParseIP("0.0.0.0")},
							{IP: net.ParseIP("::")},
						},
						Endpoint: wireguardtypes.Endpoint{
							Host: net.IP(request.Info[20 : 20+4]).String(),
							Port: binary.BigEndian.Uint16(request.Info[24 : 24+2]),
						},
						PersistentKeepalive: 15,
					},
				},
			}

			service = wireguard.NewWireGuard().
				WithConfig(wireGuardConfig)
		)

		status = status.
			WithID(request.ID).
			WithFrom(key.GetAddress().String()).
			WithTo(request.To).
			WithIFace(wireGuardConfig.Name)

		if err := status.SaveToPath(statusFilePath); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1008, err.Error()),
			)
			return
		}

		if err := service.PreUp(); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1009, err.Error()),
			)
			return
		}
		if err := service.Up(); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1010, err.Error()),
			)
			return
		}
		if err := service.PostUp(); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1011, err.Error()),
			)
			return
		}

		restutils.WriteResultToResponse(w, http.StatusOK, nil)
	}
}

func HandlerDisconnect(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			status         = types.NewStatus()
			statusFilePath = filepath.Join(ctx.Home(), types.StatusFilename)
		)

		if err := status.LoadFromPath(statusFilePath); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1001, err.Error()),
			)
			return
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
					restutils.WriteErrorToResponse(
						w, http.StatusInternalServerError,
						resttypes.NewError(1002, err.Error()),
					)
					return
				}
				if err := service.Down(); err != nil {
					restutils.WriteErrorToResponse(
						w, http.StatusInternalServerError,
						resttypes.NewError(1003, err.Error()),
					)
					return
				}
				if err := service.PostDown(); err != nil {
					restutils.WriteErrorToResponse(
						w, http.StatusInternalServerError,
						resttypes.NewError(1004, err.Error()),
					)
					return
				}
			}
		}

		if err := os.Remove(statusFilePath); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1005, err.Error()),
			)
			return
		}

		restutils.WriteResultToResponse(w, http.StatusOK, nil)
	}
}

func HandlerStatus(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			status         = types.NewStatus()
			statusFilePath = filepath.Join(ctx.Home(), types.StatusFilename)
		)

		if err := status.LoadFromPath(statusFilePath); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1001, err.Error()),
			)
			return
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
				upload, download, err := service.Transfer()
				if err != nil {
					restutils.WriteErrorToResponse(
						w, http.StatusInternalServerError,
						resttypes.NewError(1002, err.Error()),
					)
					return
				}

				restutils.WriteResultToResponse(w, http.StatusOK,
					&ResponseStatus{
						From:     status.From,
						ID:       status.ID,
						To:       status.To,
						Upload:   upload,
						Download: download,
					},
				)
				return
			}
		}

		restutils.WriteResultToResponse(w, http.StatusOK, nil)
	}
}
