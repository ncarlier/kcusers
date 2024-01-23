package getuser

import (
	"fmt"
	"log/slog"

	"github.com/ncarlier/kcusers/pkg/keycloak"
)

const unableToGetUser = "unable to get user"

func exec(client *keycloak.Client, uid string) error {
	slog.Debug("getting user...", "uid", uid)
	if err := getUser(client, uid); err != nil {
		slog.Error(unableToGetUser, "uid", uid, "error", err)
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
