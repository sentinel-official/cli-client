package keys

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/gorilla/mux"

	"github.com/sentinel-official/cli-client/context"
	keyringtypes "github.com/sentinel-official/cli-client/types/keyring"
	resttypes "github.com/sentinel-official/cli-client/types/rest"
	restutils "github.com/sentinel-official/cli-client/utils/rest"
)

func HandlerGetKey(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			config = ctx.Config()
			vars   = mux.Vars(r)
		)

		kr, err := keyring.New(
			version.Name,
			config.Keyring.Backend,
			ctx.Home(),
			strings.NewReader(""),
		)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1001, err.Error()),
			)
			return
		}

		key, err := kr.Key(vars["name"])
		if err != nil {
			if strings.Contains(err.Error(), "could not be found") {
				restutils.WriteErrorToResponse(
					w, http.StatusNotFound,
					resttypes.NewError(1002, err.Error()),
				)
				return
			}

			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1002, err.Error()),
			)
			return
		}

		item := keyringtypes.NewKeyFromRaw(key)
		restutils.WriteResultToResponse(w, http.StatusOK, item)
	}
}

func HandlerGetKeys(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			config = ctx.Config()
		)

		kr, err := keyring.New(
			version.Name,
			config.Keyring.Backend,
			ctx.Home(),
			strings.NewReader(""),
		)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1001, err.Error()),
			)
			return
		}

		list, err := kr.List()
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1002, err.Error()),
			)
			return
		}

		items := keyringtypes.NewKeysFromRaw(list)
		restutils.WriteResultToResponse(w, http.StatusOK, items)
	}
}

func HandlerAddKey(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			config = ctx.Config()
		)

		request, err := NewRequestAddKey(r)
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

		kr, err := keyring.New(
			version.Name,
			config.Keyring.Backend,
			ctx.Home(),
			strings.NewReader(""),
		)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1003, err.Error()),
			)
			return
		}

		key, _ := kr.Key(request.Name)
		if key != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusConflict,
				resttypes.NewError(1004, fmt.Sprintf("key with name %s already exists", request.Name)),
			)
			return
		}

		var (
			path          = hd.CreateHDPath(request.Type, request.Account, request.Index)
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

		key, err = kr.NewAccount(request.Name, request.Mnemonic, request.BIP39Password, path.String(), algorithm)
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

func HandlerSignBytes(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			config = ctx.Config()
			vars   = mux.Vars(r)
		)

		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1001, err.Error()),
			)
			return
		}

		defer r.Body.Close()

		kr, err := keyring.New(
			version.Name,
			config.Keyring.Backend,
			ctx.Home(),
			strings.NewReader(""),
		)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1002, err.Error()),
			)
			return
		}

		signature, pubKey, err := kr.Sign(vars["name"], bytes)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1003, err.Error()),
			)
			return
		}

		restutils.WriteResultToResponse(w, http.StatusOK,
			&ResponseSignTx{
				PubKey:    hex.EncodeToString(pubKey.Bytes()),
				Signature: signature,
			},
		)
	}
}

func HandlerDeleteKey(ctx *context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			config = ctx.Config()
			vars   = mux.Vars(r)
		)

		kr, err := keyring.New(
			version.Name,
			config.Keyring.Backend,
			ctx.Home(),
			strings.NewReader(""),
		)
		if err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1001, err.Error()),
			)
			return
		}

		if err := kr.Delete(vars["name"]); err != nil {
			restutils.WriteErrorToResponse(
				w, http.StatusInternalServerError,
				resttypes.NewError(1002, err.Error()),
			)
			return
		}

		restutils.WriteResultToResponse(w, http.StatusOK, nil)
	}
}
