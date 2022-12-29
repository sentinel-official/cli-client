package context

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/spf13/cobra"

	restrequests "github.com/sentinel-official/cli-client/rest/requests"
	restresponses "github.com/sentinel-official/cli-client/rest/responses"
	restroutes "github.com/sentinel-official/cli-client/rest/routes"
	clitypes "github.com/sentinel-official/cli-client/types"
	cliutils "github.com/sentinel-official/cli-client/utils"
)

type KeyringContext struct {
	http.Client
	Backend string
	Home    string
	URL     string
}

func NewKeyringContextFromCmd(cmd *cobra.Command) (ctx KeyringContext, err error) {
	ctx.Client, err = clitypes.GetHTTPClientFromCmd(cmd)
	if err != nil {
		return ctx, err
	}

	ctx.Backend, err = cmd.Flags().GetString(clitypes.FlagKeyringBackend)
	if err != nil {
		return ctx, err
	}

	ctx.Home, err = cmd.Flags().GetString(clitypes.FlagKeyringHome)
	if err != nil {
		return ctx, err
	}

	ctx.URL, err = cliutils.ReadLineFromFile(filepath.Join(ctx.Home, "url.txt"))
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (c KeyringContext) WithClient(v http.Client) KeyringContext {
	c.Client = v
	return c
}

func (c KeyringContext) WithBackend(v string) KeyringContext {
	c.Backend = v
	return c
}

func (c KeyringContext) WithHome(v string) KeyringContext {
	c.Home = v
	return c
}

func (c KeyringContext) WithURL(v string) KeyringContext {
	c.URL = v
	return c
}

func (c *KeyringContext) GetPasswordAndAddress(r *bufio.Reader, name string) (string, sdk.AccAddress, error) {
	password, err := cliutils.GetPassword(c.Backend, r)
	if err != nil {
		return "", nil, err
	}

	accAddr, err := c.GetAddress(password, name)
	if err != nil {
		return "", nil, err
	}

	return password, accAddr, nil
}

func (c *KeyringContext) GetAddress(password, name string) (sdk.AccAddress, error) {
	key, err := c.GetKey(password, name)
	if err != nil {
		return nil, err
	}

	accAddr, err := base64.StdEncoding.DecodeString(key.Address)
	if err != nil {
		return nil, err
	}

	return accAddr, nil
}

func (c *KeyringContext) GetKeys(password string) (clitypes.Keys, error) {
	path, err := url.JoinPath(c.URL, restroutes.GetKeys)
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(
		&restrequests.GeyKeys{
			Backend:  c.Backend,
			Password: password,
		},
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(path, jsonrpc.ContentType, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	var (
		body clitypes.RestResponseBody
		res  clitypes.Keys
	)

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if body.Error != nil {
		return nil, fmt.Errorf(body.Error.Message)
	}

	buf, err = json.Marshal(body.Result)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *KeyringContext) GetKey(password, name string) (*clitypes.Key, error) {
	path, err := url.JoinPath(c.URL, restroutes.GetKey)
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(
		&restrequests.GeyKey{
			Backend:  c.Backend,
			Password: password,
			Name:     name,
		},
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(path, jsonrpc.ContentType, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	var (
		body clitypes.RestResponseBody
		res  clitypes.Key
	)

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if body.Error != nil {
		return nil, fmt.Errorf(body.Error.Message)
	}

	buf, err = json.Marshal(body.Result)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *KeyringContext) AddKey(password, name, mnemonic, bip39Password string, coinType, account, index uint32) (*clitypes.Key, error) {
	path, err := url.JoinPath(c.URL, restroutes.AddKey)
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(
		&restrequests.AddKey{
			Backend:       c.Backend,
			Password:      password,
			Name:          name,
			Mnemonic:      mnemonic,
			CoinType:      coinType,
			Account:       account,
			Index:         index,
			BIP39Password: bip39Password,
		},
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(path, jsonrpc.ContentType, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	var (
		body clitypes.RestResponseBody
		res  clitypes.Key
	)

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if body.Error != nil {
		return nil, fmt.Errorf(body.Error.Message)
	}

	buf, err = json.Marshal(body.Result)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *KeyringContext) DeleteKey(password, name string) error {
	path, err := url.JoinPath(c.URL, restroutes.DeleteKey)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(
		&restrequests.DeleteKey{
			Backend:  c.Backend,
			Password: password,
			Name:     name,
		},
	)
	if err != nil {
		return err
	}

	resp, err := c.Post(path, jsonrpc.ContentType, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	var body clitypes.RestResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}
	if body.Error != nil {
		return fmt.Errorf(body.Error.Message)
	}

	return nil
}

func (c *KeyringContext) SignMessage(password, name string, message []byte) (*restresponses.SignMessage, error) {
	path, err := url.JoinPath(c.URL, restroutes.SignMessage)
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(
		&restrequests.SignMessage{
			Backend:  c.Backend,
			Password: password,
			Name:     name,
			Message:  message,
		},
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(path, jsonrpc.ContentType, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	var (
		body clitypes.RestResponseBody
		res  restresponses.SignMessage
	)

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if body.Error != nil {
		return nil, fmt.Errorf(body.Error.Message)
	}

	buf, err = json.Marshal(body.Result)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
