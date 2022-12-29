package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	cryptohd "github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sentinel-official/cli-client/context"
	"github.com/sentinel-official/cli-client/rest/requests"
	"github.com/sentinel-official/cli-client/rest/responses"
	clitypes "github.com/sentinel-official/cli-client/types"
	cliutils "github.com/sentinel-official/cli-client/utils"
)

func GetKey(ctx *context.ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := requests.NewGeyKey(r)
		if err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusBadRequest,
				clitypes.NewRestError(1001, err.Error()),
			)
			return
		}
		if err := req.Validate(); err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusBadRequest,
				clitypes.NewRestError(1002, err.Error()),
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
			cliutils.WriteErrorToResponseBody(
				w, http.StatusInternalServerError,
				clitypes.NewRestError(1003, err.Error()),
			)
			return
		}

		key, err := kr.Key(req.Name)
		if err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusInternalServerError,
				clitypes.NewRestError(1004, err.Error()),
			)
			return
		}

		item := clitypes.NewKeyFromRaw(key)
		cliutils.WriteResultToResponseBody(w, http.StatusOK, item)
	}
}

func GetKeys(ctx *context.ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := requests.NewGeyKeys(r)
		if err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusBadRequest,
				clitypes.NewRestError(1001, err.Error()),
			)
			return
		}
		if err := req.Validate(); err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusBadRequest,
				clitypes.NewRestError(1002, err.Error()),
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
			cliutils.WriteErrorToResponseBody(
				w, http.StatusInternalServerError,
				clitypes.NewRestError(1003, err.Error()),
			)
			return
		}

		list, err := kr.List()
		if err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusInternalServerError,
				clitypes.NewRestError(1004, err.Error()),
			)
			return
		}

		items := clitypes.NewKeysFromRaw(list)
		cliutils.WriteResultToResponseBody(w, http.StatusOK, items)
	}
}

func AddKey(ctx *context.ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := requests.NewAddKey(r)
		if err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusBadRequest,
				clitypes.NewRestError(1001, err.Error()),
			)
			return
		}
		if err := req.Validate(); err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusBadRequest,
				clitypes.NewRestError(1002, err.Error()),
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
			cliutils.WriteErrorToResponseBody(
				w, http.StatusInternalServerError,
				clitypes.NewRestError(1003, err.Error()),
			)
			return
		}

		key, _ := kr.Key(req.Name)
		if key != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusConflict,
				clitypes.NewRestError(1004, fmt.Sprintf("key with name %s already exists", req.Name)),
			)
			return
		}

		var (
			path          = cryptohd.CreateHDPath(req.CoinType, req.Account, req.Index)
			algorithms, _ = kr.SupportedAlgorithms()
		)

		algorithm, err := keyring.NewSigningAlgoFromString(string(cryptohd.Secp256k1Type), algorithms)
		if err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusInternalServerError,
				clitypes.NewRestError(1005, err.Error()),
			)
			return
		}

		key, err = kr.NewAccount(req.Name, req.Mnemonic, req.BIP39Password, path.String(), algorithm)
		if err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusInternalServerError,
				clitypes.NewRestError(1006, err.Error()),
			)
			return
		}

		item := clitypes.NewKeyFromRaw(key)
		cliutils.WriteResultToResponseBody(w, http.StatusCreated, item)
	}
}

func SignMessage(ctx *context.ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := requests.NewSignMessage(r)
		if err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusBadRequest,
				clitypes.NewRestError(1001, err.Error()),
			)
			return
		}
		if err := req.Validate(); err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusBadRequest,
				clitypes.NewRestError(1002, err.Error()),
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
			cliutils.WriteErrorToResponseBody(
				w, http.StatusInternalServerError,
				clitypes.NewRestError(1003, err.Error()),
			)
			return
		}

		signature, pubKey, err := kr.Sign(req.Name, req.Message)
		if err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusInternalServerError,
				clitypes.NewRestError(1004, err.Error()),
			)
			return
		}

		cliutils.WriteResultToResponseBody(w, http.StatusOK,
			&responses.SignMessage{
				PubKey:    base64.StdEncoding.EncodeToString(pubKey.Bytes()),
				Signature: base64.StdEncoding.EncodeToString(signature),
			},
		)
	}
}

func DeleteKey(ctx *context.ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := requests.NewDeleteKey(r)
		if err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusBadRequest,
				clitypes.NewRestError(1001, err.Error()),
			)
			return
		}
		if err := req.Validate(); err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusBadRequest,
				clitypes.NewRestError(1002, err.Error()),
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
			cliutils.WriteErrorToResponseBody(
				w, http.StatusInternalServerError,
				clitypes.NewRestError(1003, err.Error()),
			)
			return
		}

		if err := kr.Delete(req.Name); err != nil {
			cliutils.WriteErrorToResponseBody(
				w, http.StatusInternalServerError,
				clitypes.NewRestError(1004, err.Error()),
			)
			return
		}

		cliutils.WriteResultToResponseBody(w, http.StatusOK, nil)
	}
}
