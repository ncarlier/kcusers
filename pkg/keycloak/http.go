package keycloak

import (
	"net/http"
	"time"
)

var defaultHTTPClient = http.Client{
	Timeout: 5 * time.Second,
}
