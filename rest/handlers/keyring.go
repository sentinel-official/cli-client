package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sentinel-official/cli-client/context"
	restrequests "github.com/sentinel-official/cli-client/rest/requests"
	restresponses "github.com/sentinel-official/cli-client/rest/responses"
	keyringtypes "github.com/sentinel-official/cli-client/types/keyring"
	resttypes "github.com/sentinel-official/cli-client/types/rest"
	restutils "github.com/sentinel-official/cli-client/utils/rest"
)

func GetKey(ctx *context.ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := restrequests.NewGeyKey(r)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusBadRequest,
				resttypes.NewError(1001, err.Error()),
			)
			return
		}
		if err := req.Validate(); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusBadRequest,
				resttypes.NewError(1002, err.Error()),
			)
			return
		}

		kr, err := keyring.New(
			sdk.KeyringServiceName(),
			req.Backend,
			ctx.Home(),
			strings.NewReader(strings.Repeat(req.Password+"\n", 4)),
		)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1003, err.Error()),
			)
			return
		}

		key, err := kr.Key(req.Name)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1004, err.Error()),
			)
			return
		}

		item := keyringtypes.NewKeyFromRaw(key)
		restutils.WriteResultToResponse(w, http.StatusOK, item)
	}
}

func GetKeys(ctx *context.ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := restrequests.NewGeyKeys(r)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusBadRequest,
				resttypes.NewError(1001, err.Error()),
			)
			return
		}
		if err := req.Validate(); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusBadRequest,
				resttypes.NewError(1002, err.Error()),
			)
			return
		}

		kr, err := keyring.New(
			sdk.KeyringServiceName(),
			req.Backend,
			ctx.Home(),
			strings.NewReader(strings.Repeat(req.Password+"\n", 4)),
		)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1003, err.Error()),
			)
			return
		}

		list, err := kr.List()
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1004, err.Error()),
			)
			return
		}

		items := keyringtypes.NewKeysFromRaw(list)
		restutils.WriteResultToResponse(w, http.StatusOK, items)
	}
}

func AddKey(ctx *context.ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := restrequests.NewAddKey(r)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusBadRequest,
				resttypes.NewError(1001, err.Error()),
			)
			return
		}
		if err := req.Validate(); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusBadRequest,
				resttypes.NewError(1002, err.Error()),
			)
			return
		}

		kr, err := keyring.New(
			sdk.KeyringServiceName(),
			req.Backend,
			ctx.Home(),
			strings.NewReader(strings.Repeat(req.Password+"\n", 4)),
		)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1003, err.Error()),
			)
			return
		}

		key, _ := kr.Key(req.Name)
		if key != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusConflict,
				resttypes.NewError(1004, fmt.Sprintf("key with name %s already exists", req.Name)),
			)
			return
		}

		var (
			path          = hd.CreateHDPath(req.CoinType, req.Account, req.Index)
			algorithms, _ = kr.SupportedAlgorithms()
		)

		algorithm, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), algorithms)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1005, err.Error()),
			)
			return
		}

		key, err = kr.NewAccount(req.Name, req.Mnemonic, req.BIP39Password, path.String(), algorithm)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1006, err.Error()),
			)
			return
		}

		item := keyringtypes.NewKeyFromRaw(key)
		restutils.WriteResultToResponse(w, http.StatusCreated, item)
	}
}

func SignBytes(ctx *context.ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := restrequests.NewSignBytes(r)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusBadRequest,
				resttypes.NewError(1001, err.Error()),
			)
			return
		}
		if err := req.Validate(); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusBadRequest,
				resttypes.NewError(1002, err.Error()),
			)
			return
		}

		kr, err := keyring.New(
			sdk.KeyringServiceName(),
			req.Backend,
			ctx.Home(),
			strings.NewReader(strings.Repeat(req.Password+"\n", 4)),
		)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1003, err.Error()),
			)
			return
		}

		signature, pubKey, err := kr.Sign(req.Name, req.Bytes)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1004, err.Error()),
			)
			return
		}

		restutils.WriteResultToResponse(w, http.StatusOK,
			&restresponses.SignBytes{
				PubKey:    base64.StdEncoding.EncodeToString(pubKey.Bytes()),
				Signature: base64.StdEncoding.EncodeToString(signature),
			},
		)
	}
}

func DeleteKey(ctx *context.ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := restrequests.NewDeleteKey(r)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusBadRequest,
				resttypes.NewError(1001, err.Error()),
			)
			return
		}
		if err := req.Validate(); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusBadRequest,
				resttypes.NewError(1002, err.Error()),
			)
			return
		}

		kr, err := keyring.New(
			sdk.KeyringServiceName(),
			req.Backend,
			ctx.Home(),
			strings.NewReader(strings.Repeat(req.Password+"\n", 4)),
		)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1003, err.Error()),
			)
			return
		}

		if err := kr.Delete(req.Name); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1004, err.Error()),
			)
			return
		}

		restutils.WriteResultToResponse(w, http.StatusOK, nil)
	}
}
