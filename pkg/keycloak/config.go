package keycloak

// Config is the Keycloak configuration section
type Config struct {
	Authority    string
	Realm        string
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
}
