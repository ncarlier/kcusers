package resetpassword

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ncarlier/kcusers/pkg/keycloak"
)

const unableToResetPassword = "unable to reset password"

func exec(client *keycloak.Client, uid, username, password string, temporary bool) error {
	resolvedUID := uid
	if resolvedUID == "" {
		slog.Debug("looking up user by username...", "username", username)
		foundUID, err := client.LookupUIDByUsername(username)
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
