package keycloak

import (
	"time"

	"github.com/ncarlier/kcusers/pkg/toml"
)

const (
	defaultAuthority = "http://localhost:8080"
	defaultRealm     = "master"
	defaultClientID  = "test-client"
	defaultTimeout   = 5 * time.Second
	defaultCache     = ".kcusers-token.json"
)

// Config is the Keycloak configuration section
type Config struct {
	Authority    string
	Realm        string
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	Timeout      toml.Duration
	Cache        string `toml:"cache"`
	TLSInsecure  bool   `toml:"tls_insecure"`
}

// ApplyDefaults applies default values to empty fields
func (c *Config) ApplyDefaults() {
	if c.Authority == "" {
		c.Authority = defaultAuthority
	}
	if c.Realm == "" {
		c.Realm = defaultRealm
	}
	if c.ClientID == "" {
		c.ClientID = defaultClientID
	}
	if c.Timeout.Duration == 0 {
		c.Timeout.Duration = defaultTimeout
	}
	if c.Cache == "" {
		c.Cache = defaultCache
	}
}

func NewDefaultConfig() Config {
	return Config{
		Authority: defaultAuthority,
		Realm:     defaultRealm,
		ClientID:  defaultClientID,
		Timeout: toml.Duration{
			Duration: defaultTimeout,
		},
		Cache:       ".kcusers-token.json",
		TLSInsecure: false,
	}
}
