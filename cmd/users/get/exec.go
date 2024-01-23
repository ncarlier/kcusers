package getuser

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

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
	endpoint := fmt.Sprintf("%s/users/%s", client.GetAdminBaseURL(), uid)
	req, err := http.NewRequest("GET", endpoint, http.NoBody)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("invalid response: %s", res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
