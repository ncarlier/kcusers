package oidc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorResponse JSON error response
type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

func decodeErrorResponse(resp *http.Response) error {
	payload := &ErrorResponse{}
	if err := json.NewDecoder(resp.Body).Decode(payload); err != nil {
		return fmt.Errorf("invalid HTTP response code: %s", resp.Status)
	}
	return fmt.Errorf("invalid HTTP response (%s): %s", resp.Status, payload.Description)
}
