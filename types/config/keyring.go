package config

import (
	"github.com/pkg/errors"
)

type KeyringConfig struct {
	Backend string `json:"backend" mapstructure:"backend"`
}

func NewKeyringConfig() *KeyringConfig {
	return &KeyringConfig{}
}

func (c *KeyringConfig) Validate() error {
	if c.Backend == "" {
		return errors.New("backend cannot be empty")
	}
	if c.Backend != "os" && c.Backend != "test" {
		return errors.New("backend must be either os or test")
	}

	return nil
}

func (c *KeyringConfig) WithDefaultValues() *KeyringConfig {
	c.Backend = "os"

	return c
}
