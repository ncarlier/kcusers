package oidc

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
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

type OIDCClientCredentialProvider struct {
	clientID      string
	clientSecret  string
	tokenURL      *url.URL
	tokenResponse *TokenResponse
	timer         *time.Timer
	retryDelay    time.Duration
	delay         time.Duration
}

func NewOIDCClientCredentialProvider(clientID, clientSecret string, tokenURL *url.URL) *OIDCClientCredentialProvider {
	return &OIDCClientCredentialProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		tokenURL:     tokenURL,
		retryDelay:   defaultRetryDelay,
	}
}

func (o *OIDCClientCredentialProvider) GetAccessToken() string {
	return o.tokenResponse.AccessToken
}

func (o *OIDCClientCredentialProvider) fetchToken() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", o.clientID)
	data.Set("client_secret", o.clientSecret)

	req, err := http.NewRequest("POST", o.tokenURL.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := defaultHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&o.tokenResponse); err != nil {
		return err
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
	if err := o.refreshToken(); err != nil {
		return err
	}
	go o.setRefreshTokenTimer()
	return nil
}

func (o *OIDCClientCredentialProvider) Stop() {
	if o.timer != nil {
		o.timer.Stop()
	}
}
