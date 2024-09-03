package keycloak

import (
	"time"

	"github.com/ncarlier/kcusers/pkg/toml"
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

func NewDefaultConfig() Config {
	return Config{
		Authority:    "http://localhost:8080",
		Realm:        "test",
		ClientID:     "test",
		ClientSecret: "",
		Timeout: toml.Duration{
			Duration: 5 * time.Second,
		},
		Cache:       "",
		TLSInsecure: false,
	}
}
