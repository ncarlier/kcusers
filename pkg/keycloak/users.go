package keycloak

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type UserRepresentation struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Enabled  bool   `json:"enabled,omitempty"`
}

type PartialImportRepresentation struct {
	IfResourceExists string               `json:"ifResourceExists"`
	Users            []UserRepresentation `json:"users"`
}

// LookupUIDByUsername finds a user's ID by their exact username
func (c *Client) LookupUIDByUsername(username string) (string, error) {
	resource := fmt.Sprintf("/users?username=%s&exact=true", url.QueryEscape(username))
	data, err := c.AdminOperation("GET", resource)
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
		return "", fmt.Errorf("multiple users found with username: %s", username)
	}

	return users[0].ID, nil
}
