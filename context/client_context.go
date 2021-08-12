package context

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/spf13/cobra"

	restrequests "github.com/sentinel-official/cli-client/rest/requests"
	restresponses "github.com/sentinel-official/cli-client/rest/responses"
	restroutes "github.com/sentinel-official/cli-client/rest/routes"
	clitypes "github.com/sentinel-official/cli-client/types"
	keyringtypes "github.com/sentinel-official/cli-client/types/keyring"
	resttypes "github.com/sentinel-official/cli-client/types/rest"
	fileutils "github.com/sentinel-official/cli-client/utils/file"
	keyringutils "github.com/sentinel-official/cli-client/utils/keyring"
)

func PrepareHTTPClientFromCmd(cmd *cobra.Command) (c http.Client, err error) {
	timeout, err := cmd.Flags().GetDuration(clitypes.FlagTimeout)
	if err != nil {
		return c, err
	}

	return http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: timeout,
	}, nil
}

type ClientContext struct {
	http.Client
	TxContext
	KeyringBackend string
	KeyringHome    string
	KeyringURL     string
	ServiceHome    string
	ServiceURL     string
}

func NewClientContextFromCmd(cmd *cobra.Command) (ctx ClientContext, err error) {
	ctx.Client, err = PrepareHTTPClientFromCmd(cmd)
	if err != nil {
		return ctx, err
	}

	ctx.TxContext, err = NewTxContextFromCmd(cmd)
	if err != nil {
		return ctx, err
	}

	ctx.KeyringBackend, err = cmd.Flags().GetString(clitypes.FlagKeyringBackend)
	if err != nil {
		return ctx, err
	}

	ctx.KeyringHome, err = cmd.Flags().GetString(clitypes.FlagKeyringHome)
	if err != nil {
		return ctx, err
	}

	ctx.KeyringURL, err = fileutils.ReadLine(filepath.Join(ctx.KeyringHome, "url.txt"))
	if err != nil {
		return ctx, err
	}

	ctx.ServiceHome, err = cmd.Flags().GetString(clitypes.FlagServiceHome)
	if err != nil {
		return ctx, err
	}

	ctx.ServiceURL, err = fileutils.ReadLine(filepath.Join(ctx.ServiceHome, "url.txt"))
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (c ClientContext) WithHTTPClient(v http.Client) ClientContext {
	c.Client = v
	return c
}

func (c ClientContext) WithTxContext(v TxContext) ClientContext {
	c.TxContext = v
	return c
}

func (c ClientContext) WithKeyringBackend(v string) ClientContext {
	c.KeyringBackend = v
	return c
}

func (c ClientContext) WithKeyringHome(v string) ClientContext {
	c.KeyringHome = v
	return c
}

func (c ClientContext) WithKeyringURL(v string) ClientContext {
	c.KeyringURL = v
	return c
}

func (c ClientContext) WithServiceHome(v string) ClientContext {
	c.ServiceHome = v
	return c
}

func (c ClientContext) WithServiceURL(v string) ClientContext {
	c.ServiceURL = v
	return c
}

func (c *ClientContext) ReadPasswordAndGetAddress(r *bufio.Reader, name string) (string, sdk.AccAddress, error) {
	password, err := keyringutils.ReadPassword(c.KeyringBackend, r)
	if err != nil {
		return "", nil, err
	}

	address, err := c.GetAddress(password, name)
	if err != nil {
		return "", nil, err
	}

	return password, address, nil
}

func (c *ClientContext) GetAddress(password, name string) (sdk.AccAddress, error) {
	key, err := c.GetKey(password, name)
	if err != nil {
		return nil, err
	}

	address, err := base64.StdEncoding.DecodeString(key.Address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (c *ClientContext) GetKeys(password string) (keyringtypes.Keys, error) {
	var (
		resp     resttypes.Response
		result   keyringtypes.Keys
		endpoint = c.KeyringURL + restroutes.GetKeys
	)

	buf, err := json.Marshal(
		&restrequests.GeyKeys{
			Backend:  c.KeyringBackend,
			Password: password,
		},
	)
	if err != nil {
		return nil, err
	}

	r, err := c.Post(endpoint, jsonrpc.ContentType, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, fmt.Errorf(resp.Error.Message)
	}

	buf, err = json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *ClientContext) GetKey(password, name string) (*keyringtypes.Key, error) {
	var (
		resp     resttypes.Response
		result   keyringtypes.Key
		endpoint = c.KeyringURL + restroutes.GetKey
	)

	buf, err := json.Marshal(
		&restrequests.GeyKey{
			Backend:  c.KeyringBackend,
			Password: password,
			Name:     name,
		},
	)
	if err != nil {
		return nil, err
	}

	r, err := c.Post(endpoint, jsonrpc.ContentType, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, fmt.Errorf(resp.Error.Message)
	}

	buf, err = json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *ClientContext) AddKey(password, name, mnemonic, bip39Password string, coinType, account, index uint32) (*keyringtypes.Key, error) {
	var (
		resp     resttypes.Response
		result   keyringtypes.Key
		endpoint = c.KeyringURL + restroutes.AddKey
	)

	buf, err := json.Marshal(
		&restrequests.AddKey{
			Backend:       c.KeyringBackend,
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

	r, err := c.Post(endpoint, jsonrpc.ContentType, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, fmt.Errorf(resp.Error.Message)
	}

	buf, err = json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *ClientContext) DeleteKey(password, name string) error {
	var (
		resp     resttypes.Response
		endpoint = c.KeyringURL + restroutes.Delete
	)

	buf, err := json.Marshal(
		&restrequests.DeleteKey{
			Backend:  c.KeyringBackend,
			Password: password,
			Name:     name,
		},
	)
	if err != nil {
		return err
	}

	r, err := c.Post(endpoint, jsonrpc.ContentType, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return err
	}
	if resp.Error != nil {
		return fmt.Errorf(resp.Error.Message)
	}

	return nil
}

func (c *ClientContext) SignBytes(password, name string, data []byte) (*restresponses.SignBytes, error) {
	var (
		resp     resttypes.Response
		result   restresponses.SignBytes
		endpoint = c.KeyringURL + restroutes.SignBytes
	)

	buf, err := json.Marshal(
		&restrequests.SignBytes{
			Backend:  c.KeyringBackend,
			Password: password,
			Name:     name,
			Bytes:    data,
		},
	)
	if err != nil {
		return nil, err
	}

	r, err := c.Post(endpoint, jsonrpc.ContentType, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, fmt.Errorf(resp.Error.Message)
	}

	buf, err = json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *ClientContext) GetStatus() (*clitypes.Status, error) {
	var (
		resp     resttypes.Response
		result   clitypes.Status
		endpoint = c.ServiceURL + restroutes.GetStatus
	)

	r, err := c.Post(endpoint, jsonrpc.ContentType, nil)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, fmt.Errorf(resp.Error.Message)
	}

	buf, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *ClientContext) Connect(password, from, to string, id uint64, info []byte, keys [][]byte, resolvers []net.IP) error {
	var (
		resp     resttypes.Response
		endpoint = c.ServiceURL + restroutes.Connect
	)

	buf, err := json.Marshal(
		&restrequests.Connect{
			Backend:   c.KeyringBackend,
			Password:  password,
			ID:        id,
			From:      from,
			To:        to,
			Info:      info,
			Keys:      keys,
			Resolvers: resolvers,
		},
	)
	if err != nil {
		return err
	}

	r, err := c.Post(endpoint, jsonrpc.ContentType, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return err
	}
	if resp.Error != nil {
		return fmt.Errorf(resp.Error.Message)
	}

	return nil
}

func (c *ClientContext) Disconnect() error {
	var (
		resp     resttypes.Response
		endpoint = c.ServiceURL + restroutes.Disconnect
	)

	r, err := c.Post(endpoint, jsonrpc.ContentType, nil)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return err
	}
	if resp.Error != nil {
		return fmt.Errorf(resp.Error.Message)
	}

	return nil
}

func (c *ClientContext) SignAndBroadcastTx(password string, messages ...sdk.Msg) (*sdk.TxResponse, error) {
	from, err := c.GetKey(password, c.From)
	if err != nil {
		return nil, err
	}

	address, err := base64.StdEncoding.DecodeString(from.Address)
	if err != nil {
		return nil, err
	}

	pubKey, err := base64.StdEncoding.DecodeString(from.PubKey)
	if err != nil {
		return nil, err
	}

	account, err := c.QueryAccount(address)
	if err != nil {
		return nil, err
	}

	txb := c.TxConfig.NewTxBuilder()
	if err := txb.SetMsgs(messages...); err != nil {
		return nil, err
	}

	txb.SetGasLimit(c.Gas)
	txb.SetMemo(c.Memo)

	if !c.GasPrices.IsZero() {
		var (
			gas  = sdk.NewDec(int64(c.Gas))
			fees = make(sdk.Coins, len(c.GasPrices))
		)

		for i, price := range c.GasPrices {
			fee := price.Amount.Mul(gas)
			fees[i] = sdk.NewCoin(price.Denom, fee.Ceil().RoundInt())
		}

		txb.SetFeeAmount(fees)
	}

	var (
		txSignature = txsigning.SignatureV2{
			PubKey: &secp256k1.PubKey{
				Key: pubKey,
			},
			Data: &txsigning.SingleSignatureData{
				SignMode:  c.TxConfig.SignModeHandler().DefaultMode(),
				Signature: nil,
			},
			Sequence: account.GetSequence(),
		}
	)

	if err := txb.SetSignatures(txSignature); err != nil {
		return nil, err
	}

	buf, err := c.TxConfig.SignModeHandler().GetSignBytes(
		c.TxConfig.SignModeHandler().DefaultMode(),
		authsigning.SignerData{
			ChainID:       c.ChainID,
			AccountNumber: account.GetAccountNumber(),
			Sequence:      account.GetSequence(),
		},
		txb.GetTx(),
	)
	if err != nil {
		return nil, err
	}

	result, err := c.SignBytes(password, c.From, buf)
	if err != nil {
		return nil, err
	}

	signature, err := base64.StdEncoding.DecodeString(result.Signature)
	if err != nil {
		return nil, err
	}

	txSignature.Data = &txsigning.SingleSignatureData{
		SignMode:  c.TxConfig.SignModeHandler().DefaultMode(),
		Signature: signature,
	}

	if err := txb.SetSignatures(txSignature); err != nil {
		return nil, err
	}

	buf, err = c.TxConfig.TxEncoder()(txb.GetTx())
	if err != nil {
		return nil, err
	}

	res, err := c.QueryContext.BroadcastTx(buf)
	if err != nil {
		return nil, err
	}

	res.Logs = nil
	return res, nil
}
