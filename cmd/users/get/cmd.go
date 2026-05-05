package getuser

import (
	"errors"
	"flag"
	"fmt"

	"github.com/ncarlier/kcusers/cmd"
	"github.com/ncarlier/kcusers/internal/config"
	"github.com/ncarlier/kcusers/pkg/keycloak"
	"github.com/ncarlier/kcusers/pkg/uuid"
)

const cmdName = "get-user"

type GetUserCmd struct {
	uid      string
	username string
	flagSet  *flag.FlagSet
}

func (c *GetUserCmd) Exec(args []string, conf *config.Config) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}

	client, err := keycloak.NewKeycloakClient(&conf.Keycloak)
	if err != nil {
		return fmt.Errorf("unable to create Keycloak client: %w", err)
	}
	defer client.Stop()

	if c.uid == "" && c.username == "" {
		return errors.New("either uid or username must be provided")
	}

	if c.uid != "" && !uuid.IsUUID(c.uid) {
		return errors.New("invalid user ID")
	}

	return exec(client, c.uid, c.username)
}

func (c *GetUserCmd) Usage() {
	cmd.PrintCmdUsage(c.flagSet.Output(), cmdName, "Get user details")
}

func newGetUserCmd() cmd.Cmd {
	result := &GetUserCmd{}
	result.flagSet = flag.NewFlagSet(cmdName, flag.ExitOnError)
	result.flagSet.StringVar(&result.uid, "uid", "", "User ID")
	result.flagSet.StringVar(&result.username, "username", "", "Username (if uid is not provided)")
	return result
}

func init() {
	cmd.Add(cmdName, newGetUserCmd)
}
