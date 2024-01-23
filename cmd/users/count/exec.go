package countusers

import (
	"fmt"
	"log/slog"

	"github.com/ncarlier/kcusers/pkg/keycloak"
)

const unableToCountUsers = "unable to count users"

func exec(client *keycloak.Client) error {
	slog.Debug("counting users...")
	if err := countUsers(client); err != nil {
		slog.Error(unableToCountUsers, "error", err)
		return fmt.Errorf("%s: %w", unableToCountUsers, err)
	}
	return nil
}

func countUsers(client *keycloak.Client) error {
	data, err := client.AdminOperation("GET", "/users/count")
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
