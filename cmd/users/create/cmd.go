package createuser

import (
	"errors"
	"flag"
	"fmt"

	"github.com/ncarlier/kcusers/cmd"
	"github.com/ncarlier/kcusers/internal/config"
	"github.com/ncarlier/kcusers/pkg/keycloak"
	"github.com/ncarlier/kcusers/pkg/uuid"
)

const cmdName = "create-user"

type CreateUserCmd struct {
	username string
	uid      string
	flagSet  *flag.FlagSet
}

func (c *CreateUserCmd) Exec(args []string, conf *config.Config) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}

	client, err := keycloak.NewKeycloakClient(&conf.Keycloak)
	if err != nil {
		return fmt.Errorf("unable to create Keycloak client: %w", err)
	}
	defer client.Stop()

	if c.username == "" {
		return errors.New("username is required")
	}

	if c.uid != "" && !uuid.IsUUID(c.uid) {
		return errors.New("invalid user ID")
	}

	return exec(client, c.username, c.uid)
}

func (c *CreateUserCmd) Usage() {
	cmd.PrintCmdUsage(c.flagSet.Output(), cmdName, "Create a new user")
}

func newCreateUserCmd() cmd.Cmd {
	result := &CreateUserCmd{}
	result.flagSet = flag.NewFlagSet(cmdName, flag.ExitOnError)
	result.flagSet.StringVar(&result.username, "username", "", "Username (required)")
	result.flagSet.StringVar(&result.uid, "uid", "", "User ID (optional, requires Keycloak v23+ for partial import)")
	return result
}

func init() {
	cmd.Add(cmdName, newCreateUserCmd)
}
