package keycloak

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/ncarlier/kcusers/pkg/oidc"
)

// Client structure
type Client struct {
	Authority     string
	Realm         string
	TokenProvider *oidc.OIDCClientCredentialProvider
	httpClient    *http.Client
}

// NewKeycloakClient creat new Keycloak client
func NewKeycloakClient(conf *Config) (*Client, error) {
	tokenEndpoint, err := url.Parse(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", conf.Authority, conf.Realm))
	if err != nil {
		return nil, fmt.Errorf("invalid Keycloak client configuration: %w", err)
	}

	httpClient := newHTTPClient(conf)

	tokenProvider, err := oidc.NewOIDCClientCredentialProvider(&oidc.OIDCClientCredentialConfig{
		TokenEndpoint: tokenEndpoint.String(),
		ClientID:      conf.ClientID,
		ClientSecret:  conf.ClientSecret,
		TokenCache:    conf.Cache,
		HttpClient:    httpClient,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to init OIDC client provider: %w", err)
	}

	if err := tokenProvider.Start(); err != nil {
		return nil, fmt.Errorf("unable to start Keycloak client token service: %w", err)
	}

	return &Client{
		Authority:     conf.Authority,
		Realm:         conf.Realm,
		TokenProvider: tokenProvider,
		httpClient:    httpClient,
	}, nil
}

// AdminOperation do HTTP operation on an resource of Keycloak Admin API
func (c *Client) AdminOperation(method, resource string) ([]byte, error) {
	req, err := http.NewRequest(method, c.GetAdminBaseURL(resource), http.NoBody)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("invalid response: %s", res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Do HTTP request with access token
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	token := c.TokenProvider.GetAccessToken()
	req.Header.Set("Authorization", "Bearer "+token)
	return c.httpClient.Do(req)
}

// GetAdminBaseURL return admin API base URL
func (c *Client) GetAdminBaseURL(resource string) string {
	path := filepath.Join("/admin/realms", c.Realm, resource)
	return c.Authority + path
}

// Stop Keycloak client
func (c *Client) Stop() {
	c.TokenProvider.Stop()
}
