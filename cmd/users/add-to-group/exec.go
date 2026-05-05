package addtogroup

import (
	"fmt"
	"log/slog"

	"github.com/ncarlier/kcusers/pkg/keycloak"
)

const unableToAddToGroup = "unable to add user to group"

func exec(client *keycloak.Client, username, groupName string) error {
	slog.Debug("looking up user by username...", "username", username)
	uid, err := client.LookupUIDByUsername(username)
	if err != nil {
		slog.Error(unableToAddToGroup, "username", username, "error", err)
		return fmt.Errorf("%s: %w", unableToAddToGroup, err)
	}

	slog.Debug("looking up group by name...", "group", groupName)
	gid, err := client.LookupGroupIDByName(groupName)
	if err != nil {
		slog.Error(unableToAddToGroup, "group", groupName, "error", err)
		return fmt.Errorf("%s: %w", unableToAddToGroup, err)
	}

	slog.Debug("adding user to group...", "uid", uid, "gid", gid)
	resource := fmt.Sprintf("/users/%s/groups/%s", uid, gid)
	if _, err := client.AdminOperation("PUT", resource); err != nil {
		slog.Error(unableToAddToGroup, "uid", uid, "gid", gid, "error", err)
		return fmt.Errorf("%s: %w", unableToAddToGroup, err)
	}

	slog.Info("user added to group successfully", "username", username, "group", groupName)
	return nil
}
