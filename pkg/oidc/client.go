package oidc

import (
	"encoding/json"
	"errors"
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

type TokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type OIDCClientCredentialConfig struct {
	TokenEndpoint string
	TokenCache    string
	ClientID      string
	ClientSecret  string
	HttpClient    *http.Client
}

type OIDCClientCredentialProvider struct {
	clientID         string
	clientSecret     string
	tokenEndpoint    string
	httpClient       *http.Client
	tokenResponse    *TokenResponse
	tokenFile        *os.File
	tokenFileCreated bool
	timer            *time.Timer
	retryDelay       time.Duration
	delay            time.Duration
}

func NewOIDCClientCredentialProvider(cfg *OIDCClientCredentialConfig) (*OIDCClientCredentialProvider, error) {
	provider := &OIDCClientCredentialProvider{
		clientID:      cfg.ClientID,
		clientSecret:  cfg.ClientSecret,
		tokenEndpoint: cfg.TokenEndpoint,
		retryDelay:    defaultRetryDelay,
		httpClient:    cfg.HttpClient,
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

func (o *OIDCClientCredentialProvider) GetAccessToken() string {
	return o.tokenResponse.AccessToken
}

func (o *OIDCClientCredentialProvider) fetchToken() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", o.clientID)
	data.Set("client_secret", o.clientSecret)

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
	if o.tokenFile != nil {
		slog.Debug("saving token for future usage...", "file", o.tokenFile.Name())
		if err := o.tokenFile.Truncate(0); err != nil {
			return err
		}
		if err := json.NewEncoder(o.tokenFile).Encode(&o.tokenResponse); err != nil {
			return err
		}
	}
	return nil
}

func (o *OIDCClientCredentialProvider) loadToken() error {
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
			if time.Now().After(expiresAt) {
				slog.Debug("token file expired", "file", o.tokenFile.Name(), "expires_at", expiresAt)
				o.tokenResponse = nil
			} else {
				o.delay = time.Until(expiresAt)
				slog.Debug("token loaded", "file", o.tokenFile.Name(), "expires_in", o.delay)
			}
		}
	}
	return nil
}

func (o *OIDCClientCredentialProvider) refreshToken() error {
	slog.Debug("refreshing token...")
	if err := o.fetchToken(); err != nil {
		slog.Error("unable to refresh token", "error", err)
		o.delay = o.retryDelay
		return err
	}
	slog.Debug("token refreshed", "expires_in", o.tokenResponse.ExpiresIn)
	expiresIn := time.Duration(o.tokenResponse.ExpiresIn) * time.Second
	o.delay = expiresIn - safeRefreshLeap
	return nil
}

func (o *OIDCClientCredentialProvider) setRefreshTokenTimer() {
	slog.Debug("refreshing token in...", "delay", o.delay)
	next := time.Now().Add(o.delay)
	o.timer = time.NewTimer(time.Until(next))
	defer o.timer.Stop()
	<-o.timer.C
	o.refreshToken()
	go o.setRefreshTokenTimer()
}

func (o *OIDCClientCredentialProvider) Start() error {
	if err := o.loadToken(); err != nil {
		return err
	}
	if o.tokenResponse == nil {
		if err := o.refreshToken(); err != nil {
			return err
		}
	}
	go o.setRefreshTokenTimer()
	return nil
}

func (o *OIDCClientCredentialProvider) Stop() {
	if o.timer != nil {
		o.timer.Stop()
	}
	if o.tokenFile != nil {
		o.tokenFile.Close()
	}
}
