package context

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/spf13/cobra"

	restrequests "github.com/sentinel-official/cli-client/rest/requests"
	restroutes "github.com/sentinel-official/cli-client/rest/routes"
	clitypes "github.com/sentinel-official/cli-client/types"
	cliutils "github.com/sentinel-official/cli-client/utils"
)

type ServiceContext struct {
	http.Client
	Home string
	URL  string
}

func NewServiceContextFromCmd(cmd *cobra.Command) (ctx ServiceContext, err error) {
	ctx.Client, err = clitypes.GetHTTPClientFromCmd(cmd)
	if err != nil {
		return ctx, err
	}

	ctx.Home, err = cmd.Flags().GetString(clitypes.FlagServiceHome)
	if err != nil {
		return ctx, err
	}

	ctx.URL, err = cliutils.ReadLineFromFile(filepath.Join(ctx.Home, "url.txt"))
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (c ServiceContext) WithHome(v string) ServiceContext {
	c.Home = v
	return c
}

func (c ServiceContext) WithURL(v string) ServiceContext {
	c.URL = v
	return c
}

func (c *ServiceContext) GetStatus() (*clitypes.ServiceStatus, error) {
	path, err := url.JoinPath(c.URL, restroutes.GetStatus)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(path, jsonrpc.ContentType, nil)
	if err != nil {
		return nil, err
	}

	var (
		body clitypes.RestResponseBody
		res  clitypes.ServiceStatus
	)

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if body.Error != nil {
		return nil, fmt.Errorf(body.Error.Message)
	}

	buf, err := json.Marshal(body.Result)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *ServiceContext) Connect(backend, password, from, to string, id uint64, info []byte, keys [][]byte, resolvers []net.IP) error {
	path, err := url.JoinPath(c.URL, restroutes.Connect)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(
		&restrequests.Connect{
			Backend:   backend,
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

func (c *ServiceContext) Disconnect() error {
	path, err := url.JoinPath(c.URL, restroutes.Disconnect)
	if err != nil {
		return err
	}

	resp, err := c.Post(path, jsonrpc.ContentType, nil)
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
