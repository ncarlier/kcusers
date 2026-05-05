package addtogroup

import (
	"errors"
	"flag"
	"fmt"

	"github.com/ncarlier/kcusers/cmd"
	"github.com/ncarlier/kcusers/internal/config"
	"github.com/ncarlier/kcusers/pkg/keycloak"
)

const cmdName = "add-to-group"

type AddToGroupCmd struct {
	username  string
	groupName string
	flagSet   *flag.FlagSet
}

func (c *AddToGroupCmd) Exec(args []string, conf *config.Config) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}

	if c.username == "" {
		return errors.New("username is required")
	}

	if c.groupName == "" {
		return errors.New("group is required")
	}

	client, err := keycloak.NewKeycloakClient(&conf.Keycloak)
	if err != nil {
		return fmt.Errorf("unable to create Keycloak client: %w", err)
	}
	defer client.Stop()

	return exec(client, c.username, c.groupName)
}

func (c *AddToGroupCmd) Usage() {
	cmd.PrintCmdUsage(c.flagSet.Output(), cmdName, "Add a user to a group")
}

func newAddToGroupCmd() cmd.Cmd {
	result := &AddToGroupCmd{}
	result.flagSet = flag.NewFlagSet(cmdName, flag.ExitOnError)
	result.flagSet.StringVar(&result.username, "username", "", "Username (required)")
	result.flagSet.StringVar(&result.groupName, "group", "", "Group exact name (required)")
	return result
}

func init() {
	cmd.Add(cmdName, newAddToGroupCmd)
}
