package resetpassword

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/ncarlier/kcusers/pkg/keycloak"
)

const unableToResetPassword = "unable to reset password"

func exec(client *keycloak.Client, uid, username, password string, temporary bool) error {
	resolvedUID := uid
	if resolvedUID == "" {
		slog.Debug("looking up user by username...", "username", username)
		foundUID, err := lookupUIDByUsername(client, username)
		if err != nil {
			return fmt.Errorf("%s: %w", unableToResetPassword, err)
		}
		resolvedUID = foundUID
	}

	slog.Debug("resetting user password...", "uid", resolvedUID, "temporary", temporary)
	if err := resetPassword(client, resolvedUID, password, temporary); err != nil {
		slog.Error(unableToResetPassword, "uid", resolvedUID, "error", err)
		return fmt.Errorf("%s: %w", unableToResetPassword, err)
	}

	slog.Info("password reset successfully", "uid", resolvedUID)
	return nil
}

type UserRepresentation struct {
	ID string `json:"id"`
}

func lookupUIDByUsername(client *keycloak.Client, username string) (string, error) {
	resource := fmt.Sprintf("/users?username=%s&exact=true", url.QueryEscape(username))
	data, err := client.AdminOperation("GET", resource)
	if err != nil {
		return "", err
	}

	var users []UserRepresentation
	if err := json.Unmarshal(data, &users); err != nil {
		return "", fmt.Errorf("unable to parse users response: %w", err)
	}

	if len(users) == 0 {
		return "", fmt.Errorf("user not found with username: %s", username)
	}

	if len(users) > 1 {
		// Should not happen with exact=true, but good to check
		return "", fmt.Errorf("multiple users found with username: %s", username)
	}

	return users[0].ID, nil
}

type CredentialRepresentation struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

func resetPassword(client *keycloak.Client, uid, password string, temporary bool) error {
	resource := fmt.Sprintf("/users/%s/reset-password", uid)

	payload := CredentialRepresentation{
		Type:      "password",
		Value:     password,
		Temporary: temporary,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = client.AdminOperationWithBody("PUT", resource, bytes.NewBuffer(payloadBytes))
	return err
}
