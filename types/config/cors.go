package config

type CORSConfig struct {
	AllowedOrigins string `json:"allowed_origins" mapstructure:"allowed_origins"`
}

func NewCORSConfig() *CORSConfig {
	return &CORSConfig{}
}

func (c *CORSConfig) Validate() error {
	return nil
}

func (c *CORSConfig) WithDefaultValues() *CORSConfig {
	c.AllowedOrigins = "*"

	return c
}
