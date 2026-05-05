package createuser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ncarlier/kcusers/pkg/keycloak"
)

const unableToCreateUser = "unable to create user"

func exec(client *keycloak.Client, username, uid string) error {
	slog.Debug("creating user...", "username", username, "uid", uid)

	if uid != "" {
		if err := createUserWithUUID(client, username, uid); err != nil {
			slog.Error(unableToCreateUser, "username", username, "uid", uid, "error", err)
			return fmt.Errorf("%s: %w", unableToCreateUser, err)
		}
	} else {
		if err := createUser(client, username); err != nil {
			slog.Error(unableToCreateUser, "username", username, "error", err)
			return fmt.Errorf("%s: %w", unableToCreateUser, err)
		}
	}

	slog.Info("user created", "username", username, "uid", uid)
	return nil
}

func createUser(client *keycloak.Client, username string) error {
	resource := "/users"

	user := keycloak.UserRepresentation{
		Username: username,
		Enabled:  true,
	}

	payload, err := json.Marshal(user)
	if err != nil {
		return err
	}

	data, err := client.AdminOperationWithBody("POST", resource, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	if len(data) > 0 {
		fmt.Println(string(data))
	}

	return nil
}

func createUserWithUUID(client *keycloak.Client, username, uid string) error {
	resource := "/partialImport"

	importPayload := keycloak.PartialImportRepresentation{
		IfResourceExists: "FAIL",
		Users: []keycloak.UserRepresentation{
			{
				ID:       uid,
				Username: username,
				Enabled:  true,
			},
		},
	}

	payload, err := json.Marshal(importPayload)
	if err != nil {
		return err
	}

	data, err := client.AdminOperationWithBody("POST", resource, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	if len(data) > 0 {
		fmt.Println(string(data))
	}

	return nil
}
