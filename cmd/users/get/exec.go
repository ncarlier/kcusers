package getuser

import (
	"fmt"
	"log/slog"

	"github.com/ncarlier/kcusers/pkg/keycloak"
)

const unableToGetUser = "unable to get user"

func exec(client *keycloak.Client, uid, username string) error {
	resolvedUID := uid
	if resolvedUID == "" {
		slog.Debug("looking up user by username...", "username", username)
		foundUID, err := client.LookupUIDByUsername(username)
		if err != nil {
			return fmt.Errorf("%s: %w", unableToGetUser, err)
		}
		resolvedUID = foundUID
	}

	slog.Debug("getting user...", "uid", resolvedUID)
	if err := getUser(client, resolvedUID); err != nil {
		slog.Error(unableToGetUser, "uid", resolvedUID, "error", err)
		return fmt.Errorf("%s: %w", unableToGetUser, err)
	}
	return nil
}

func getUser(client *keycloak.Client, uid string) error {
	resource := fmt.Sprintf("/users/%s", uid)
	data, err := client.AdminOperation("GET", resource)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
