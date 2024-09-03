package keycloak

import (
	"crypto/tls"
	"net/http"
)

func newHTTPClient(config *Config) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.TLSInsecure},
		Proxy:           http.ProxyFromEnvironment,
	}
	return &http.Client{
		Timeout:   config.Timeout.Duration,
		Transport: tr,
	}
}
