package keycloak

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ncarlier/kcusers/pkg/oidc"
)

// Client structure
type Client struct {
	Authority     string
	Realm         string
	TokenProvider *oidc.OIDCClientCredentialProvider
}

// NewKeycloakClient creat new Keycloak client
func NewKeycloakClient(conf *Config) (*Client, error) {
	tokenEndpoint, err := url.Parse(fmt.Sprintf("%s/auth/realms/%s/protocol/openid-connect/token", conf.Authority, conf.Realm))
	if err != nil {
		return nil, fmt.Errorf("invalid Keycloak client configuration: %w", err)
	}
	tokenProvider := oidc.NewOIDCClientCredentialProvider(
		conf.ClientID,
		conf.ClientSecret,
		tokenEndpoint,
	)
	if err := tokenProvider.Start(); err != nil {
		return nil, fmt.Errorf("unable to start Keycloak client token service: %w", err)
	}
	client := &Client{
		Authority:     conf.Authority,
		Realm:         conf.Realm,
		TokenProvider: tokenProvider,
	}
	return client, nil
}

// Do HTTP request with access token
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	token := c.TokenProvider.GetAccessToken()
	req.Header.Set("Authorization", "Bearer "+token)
	return defaultHTTPClient.Do(req)
}

// GetAdminBaseURL return admin API base URL
func (c *Client) GetAdminBaseURL() string {
	return fmt.Sprintf("%s/auth/admin/realms/%s", c.Authority, c.Realm)
}

// Stop Keycloak client
func (c *Client) Stop() {
	c.TokenProvider.Stop()
}
