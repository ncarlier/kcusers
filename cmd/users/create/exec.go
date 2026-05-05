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

type UserRepresentation struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username"`
	Enabled  bool   `json:"enabled"`
}

type PartialImportRepresentation struct {
	IfResourceExists string               `json:"ifResourceExists"`
	Users            []UserRepresentation `json:"users"`
}

func createUser(client *keycloak.Client, username string) error {
	resource := "/users"

	user := UserRepresentation{
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

	importPayload := PartialImportRepresentation{
		IfResourceExists: "FAIL",
		Users: []UserRepresentation{
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
