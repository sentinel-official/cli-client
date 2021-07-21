package config

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	ct = strings.TrimSpace(`
listen_on = "{{ .ListenOn }}"
token = "{{ .Token }}"
version = {{ .Version }}

[cors]
allowed_origins = "{{ .CORS.AllowedOrigins }}"

[keyring]
backend = "{{ .Keyring.Backend }}"
	`)

	t = func() *template.Template {
		t, err := template.New("").Parse(ct)
		if err != nil {
			panic(err)
		}

		return t
	}()
)

type Config struct {
	ListenOn string         `json:"listen_on" mapstructure:"listen_on"`
	Token    string         `json:"token" mapstructure:"token"`
	Version  int64          `json:"version" mapstructure:"version"`
	CORS     *CORSConfig    `json:"cors" mapstructure:"cors"`
	Keyring  *KeyringConfig `json:"keyring" mapstructure:"keyring"`
}

func NewConfig() *Config {
	return &Config{
		CORS:    NewCORSConfig(),
		Keyring: NewKeyringConfig(),
	}
}

func (c *Config) WithDefaultValues() *Config {
	c.ListenOn = "127.0.0.1:9090"
	c.Token = ""
	c.Version = 1

	c.CORS = c.CORS.WithDefaultValues()
	c.Keyring = c.Keyring.WithDefaultValues()

	return c
}

func (c *Config) SaveToPath(path string) error {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		return err
	}

	return ioutil.WriteFile(path, buffer.Bytes(), 0600)
}

func (c *Config) String() string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}

	return buffer.String()
}

func (c *Config) Validate() error {
	if c.ListenOn == "" {
		return errors.New("listen_on cannot be empty")
	}
	if err := c.CORS.Validate(); err != nil {
		return errors.Wrapf(err, "invalid section cors")
	}
	if err := c.Keyring.Validate(); err != nil {
		return errors.Wrapf(err, "invalid section keyring")
	}

	return nil
}

func ReadInConfig(v *viper.Viper) (*Config, error) {
	cfg := NewConfig().WithDefaultValues()
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
