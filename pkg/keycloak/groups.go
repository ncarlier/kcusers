package keycloak

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type GroupRepresentation struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// LookupGroupIDByName searches for a Keycloak group by its exact name and returns its ID
func (c *Client) LookupGroupIDByName(name string) (string, error) {
	resource := fmt.Sprintf("/groups?search=%s&exact=true", url.QueryEscape(name))

	data, err := c.AdminOperation("GET", resource)
	if err != nil {
		return "", fmt.Errorf("failed to search for group: %w", err)
	}

	var groups []GroupRepresentation
	if err := json.Unmarshal(data, &groups); err != nil {
		return "", fmt.Errorf("failed to decode groups response: %w", err)
	}

	for _, group := range groups {
		if group.Name == name {
			return group.ID, nil
		}
	}

	return "", fmt.Errorf("group %q not found", name)
}
