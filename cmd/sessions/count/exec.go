package countsessions

import (
	"fmt"
	"log/slog"

	"github.com/ncarlier/kcusers/pkg/keycloak"
)

const unableToCountSessions = "unable to count sessions"

func exec(client *keycloak.Client, cid string) error {
	slog.Debug("getting sessions count...", "cid", cid)
	if err := countSessions(client, cid); err != nil {
		slog.Error(unableToCountSessions, "cid", cid, "error", err)
		return fmt.Errorf("%s: %w", unableToCountSessions, err)
	}
	return nil
}

func countSessions(client *keycloak.Client, cid string) error {
	resource := fmt.Sprintf("/clients/%s/session-count", cid)
	data, err := client.AdminOperation("GET", resource)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
