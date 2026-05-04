package oidc

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultRetryDelay = 5 * time.Second
const safeRefreshLeap = 5 * time.Second

func generateCodeVerifier() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func generateCodeChallenge(verifier string) string {
	h := sha256.New()
	h.Write([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

type TokenResponse struct {
	TokenType        string `json:"token_type"`
	AccessToken      string `json:"access_token"`
	IDToken          string `json:"id_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	RefreshExpiresIn int    `json:"refresh_expires_in,omitempty"`
}

type OIDCDeviceCodeConfig struct {
	DeviceAuthEndpoint string
	TokenEndpoint      string
	TokenCache         string
	ClientID           string
	ClientSecret       string
	HttpClient         *http.Client
}

type DeviceAuthResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

type OIDCDeviceCodeProvider struct {
	clientID           string
	clientSecret       string
	deviceAuthEndpoint string
	tokenEndpoint      string
	httpClient         *http.Client
	tokenResponse      *TokenResponse
	tokenFile          *os.File
	tokenFileCreated   bool
	timer              *time.Timer
	retryDelay         time.Duration
	delay              time.Duration
	codeVerifier       string
}

func NewOIDCDeviceCodeProvider(cfg *OIDCDeviceCodeConfig) (*OIDCDeviceCodeProvider, error) {
	provider := &OIDCDeviceCodeProvider{
		clientID:           cfg.ClientID,
		clientSecret:       cfg.ClientSecret,
		deviceAuthEndpoint: cfg.DeviceAuthEndpoint,
		tokenEndpoint:      cfg.TokenEndpoint,
		retryDelay:         defaultRetryDelay,
		httpClient:         cfg.HttpClient,
	}
	if provider.httpClient == nil {
		provider.httpClient = http.DefaultClient
	}

	if cfg.TokenCache != "" {
		if _, err := os.Stat(cfg.TokenCache); err == nil {
			if provider.tokenFile, err = os.OpenFile(cfg.TokenCache, os.O_RDWR, 0o644); err != nil {
				return nil, err
			}
			provider.tokenFileCreated = false
		} else if errors.Is(err, os.ErrNotExist) {
			if provider.tokenFile, err = os.Create(cfg.TokenCache); err != nil {
				return nil, err
			}
			provider.tokenFileCreated = true
		} else {
			return nil, err
		}
	}

	return provider, nil
}

func (o *OIDCDeviceCodeProvider) GetAccessToken() string {
	return o.tokenResponse.AccessToken
}

func (o *OIDCDeviceCodeProvider) initiateDeviceAuth() (*DeviceAuthResponse, error) {
	verifier, err := generateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("unable to generate code verifier: %w", err)
	}
	o.codeVerifier = verifier
	challenge := generateCodeChallenge(verifier)

	data := url.Values{}
	data.Set("client_id", o.clientID)
	if o.clientSecret != "" {
		data.Set("client_secret", o.clientSecret)
	}
	data.Set("code_challenge", challenge)
	data.Set("code_challenge_method", "S256")

	req, err := http.NewRequest("POST", o.deviceAuthEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, decodeErrorResponse(res)
	}

	var authRes DeviceAuthResponse
	if err := json.NewDecoder(res.Body).Decode(&authRes); err != nil {
		return nil, err
	}
	return &authRes, nil
}

func (o *OIDCDeviceCodeProvider) pollForToken(deviceCode string, interval int) error {
	data := url.Values{}
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	data.Set("client_id", o.clientID)
	if o.clientSecret != "" {
		data.Set("client_secret", o.clientSecret)
	}
	data.Set("device_code", deviceCode)
	if o.codeVerifier != "" {
		data.Set("code_verifier", o.codeVerifier)
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		req, err := http.NewRequest("POST", o.tokenEndpoint, strings.NewReader(data.Encode()))
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

		res, err := o.httpClient.Do(req)
		if err != nil {
			return err
		}

		if res.StatusCode == http.StatusOK {
			if err := json.NewDecoder(res.Body).Decode(&o.tokenResponse); err != nil {
				res.Body.Close()
				return err
			}
			res.Body.Close()
			return o.saveToken()
		}

		var oidcErr ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&oidcErr); err != nil {
			res.Body.Close()
			return err
		}
		res.Body.Close()

		if oidcErr.Error != "authorization_pending" && oidcErr.Error != "slow_down" {
			return fmt.Errorf("OIDC Error: %s - %s", oidcErr.Error, oidcErr.Description)
		}
		// If authorization_pending or slow_down, continue polling
	}
	return nil
}

func (o *OIDCDeviceCodeProvider) fetchToken() error {
	// If we have a refresh token, try to refresh first
	if o.tokenResponse != nil && o.tokenResponse.RefreshToken != "" {
		if err := o.refreshAccessToken(); err == nil {
			return nil
		}
		slog.Info("Session expired or refresh failed, starting new device flow")
	}

	authRes, err := o.initiateDeviceAuth()
	if err != nil {
		return err
	}

	fmt.Printf("\nAction required! Please visit the following URL to authenticate:\n\n%s\n\n", authRes.VerificationURIComplete)

	interval := authRes.Interval
	if interval == 0 {
		interval = 5 // default to 5 seconds if not provided
	}

	return o.pollForToken(authRes.DeviceCode, interval)
}

func (o *OIDCDeviceCodeProvider) refreshAccessToken() error {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", o.clientID)
	if o.clientSecret != "" {
		data.Set("client_secret", o.clientSecret)
	}
	data.Set("refresh_token", o.tokenResponse.RefreshToken)

	req, err := http.NewRequest("POST", o.tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := o.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return decodeErrorResponse(res)
	}

	if err := json.NewDecoder(res.Body).Decode(&o.tokenResponse); err != nil {
		return err
	}
	return o.saveToken()
}

func (o *OIDCDeviceCodeProvider) saveToken() error {
	if o.tokenFile != nil {
		slog.Debug("saving token for future usage...", "file", o.tokenFile.Name())
		if err := o.tokenFile.Truncate(0); err != nil {
			return err
		}
		if _, err := o.tokenFile.Seek(0, 0); err != nil {
			return err
		}
		if err := json.NewEncoder(o.tokenFile).Encode(&o.tokenResponse); err != nil {
			return err
		}
	}
	return nil
}

func (o *OIDCDeviceCodeProvider) loadToken() error {
	if o.tokenFile != nil && !o.tokenFileCreated {
		slog.Debug("loading token from file...", "file", o.tokenFile.Name())
		if err := json.NewDecoder(o.tokenFile).Decode(&o.tokenResponse); err != nil {
			return err
		}
		if stat, err := o.tokenFile.Stat(); err != nil {
			return err
		} else {
			expiresIn := time.Duration(o.tokenResponse.ExpiresIn) * time.Second
			expiresAt := stat.ModTime().Add(expiresIn - safeRefreshLeap)

			// If refresh token exists, calculate its expiration
			var refreshExpiresAt time.Time
			if o.tokenResponse.RefreshExpiresIn > 0 {
				refreshExpiresIn := time.Duration(o.tokenResponse.RefreshExpiresIn) * time.Second
				refreshExpiresAt = stat.ModTime().Add(refreshExpiresIn - safeRefreshLeap)
			}

			if time.Now().After(expiresAt) {
				slog.Debug("access token expired", "file", o.tokenFile.Name(), "expires_at", expiresAt)
				// If access token is expired, but refresh token is still valid (or no expiration provided)
				if o.tokenResponse.RefreshToken != "" && (refreshExpiresAt.IsZero() || time.Now().Before(refreshExpiresAt)) {
					slog.Debug("attempting to use refresh token")
					// Try to refresh token right away instead of fully zeroing it out
					if err := o.refreshAccessToken(); err != nil {
						slog.Error("failed to refresh token", "error", err)
						o.tokenResponse = nil
					} else {
						// Set the new delay
						newExpiresIn := time.Duration(o.tokenResponse.ExpiresIn) * time.Second
						o.delay = newExpiresIn - safeRefreshLeap
						slog.Debug("token refreshed from cache", "expires_in", o.delay)
					}
				} else {
					slog.Debug("refresh token expired or not available, need to re-authenticate")
					o.tokenResponse = nil
				}
			} else {
				o.delay = time.Until(expiresAt)
				slog.Debug("token loaded", "file", o.tokenFile.Name(), "expires_in", o.delay)
			}
		}
	}
	return nil
}

func (o *OIDCDeviceCodeProvider) refreshTokenTask() error {
	slog.Debug("refreshing token...")
	if err := o.refreshAccessToken(); err != nil {
		slog.Error("unable to refresh token", "error", err)
		o.delay = o.retryDelay
		return err
	}
	slog.Debug("token refreshed", "expires_in", o.tokenResponse.ExpiresIn)
	expiresIn := time.Duration(o.tokenResponse.ExpiresIn) * time.Second
	o.delay = expiresIn - safeRefreshLeap
	return nil
}

func (o *OIDCDeviceCodeProvider) setRefreshTokenTimer() {
	slog.Debug("refreshing token in...", "delay", o.delay)
	next := time.Now().Add(o.delay)
	o.timer = time.NewTimer(time.Until(next))
	defer o.timer.Stop()
	<-o.timer.C
	o.refreshTokenTask()
	go o.setRefreshTokenTimer()
}

func (o *OIDCDeviceCodeProvider) Start() error {
	if err := o.loadToken(); err != nil {
		return err
	}
	if o.tokenResponse == nil {
		if err := o.fetchToken(); err != nil {
			return err
		}
		// calculate delay after fetching fresh token
		expiresIn := time.Duration(o.tokenResponse.ExpiresIn) * time.Second
		o.delay = expiresIn - safeRefreshLeap
	}
	go o.setRefreshTokenTimer()
	return nil
}

func (o *OIDCDeviceCodeProvider) Stop() {
	if o.timer != nil {
		o.timer.Stop()
	}
	if o.tokenFile != nil {
		o.tokenFile.Close()
	}
}
