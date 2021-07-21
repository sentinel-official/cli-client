package context

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"

	configtypes "github.com/sentinel-official/cli-client/types/config"
)

type Context struct {
	config *configtypes.Config
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) WithConfig(v *configtypes.Config) *Context { c.config = v; return c }

func (c *Context) Config() *configtypes.Config { return c.config }

func (c *Context) Home() string { return viper.GetString(flags.FlagHome) }
