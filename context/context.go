package context

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	hubparams "github.com/sentinel-official/hub/params"
	"github.com/spf13/viper"

	configtypes "github.com/sentinel-official/cli-client/types/config"
)

type Context struct {
	config   *configtypes.Config
	encoding *hubparams.EncodingConfig
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) WithConfig(v *configtypes.Config) *Context         { c.config = v; return c }
func (c *Context) WithEncoding(v *hubparams.EncodingConfig) *Context { c.encoding = v; return c }

func (c *Context) Config() *configtypes.Config         { return c.config }
func (c *Context) Encoding() *hubparams.EncodingConfig { return c.encoding }

func (c *Context) Home() string { return viper.GetString(flags.FlagHome) }
