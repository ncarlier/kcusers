package config

import (
	"fmt"
	"io"
	"os"

	"github.com/influxdata/toml"
	"github.com/ncarlier/kcusers/pkg/keycloak"
)

// Config is the root of the configuration
type Config struct {
	Log      LogConfig
	Keycloak keycloak.Config
}

// LogConfig is the proxy configuration section
type LogConfig struct {
	Level  string
	Format string
}

// NewConfig create new configuration
func NewConfig() *Config {
	c := &Config{
		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
		Keycloak: keycloak.NewDefaultConfig(),
	}
	return c
}

// LoadConfig loads the given config file and applies it to c
func (c *Config) LoadConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	data = []byte(os.ExpandEnv(string(data)))
	tbl, err := toml.Parse(data)
	if err != nil {
		return err
	}

	if err = toml.UnmarshalTable(tbl, &c); err != nil {
		return fmt.Errorf("error parsing configuration: %w", err)
	}

	return nil
}
