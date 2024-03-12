package countusers

import (
	"flag"
	"fmt"

	"github.com/ncarlier/kcusers/cmd"
	"github.com/ncarlier/kcusers/internal/config"
	"github.com/ncarlier/kcusers/pkg/keycloak"
)

const cmdName = "count-users"

type CountUsersCmd struct {
	flagSet *flag.FlagSet
}

func (c *CountUsersCmd) Exec(args []string, conf *config.Config) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}

	client, err := keycloak.NewKeycloakClient(&conf.Keycloak)
	if err != nil {
		return fmt.Errorf("unable to create Keycloak client: %w", err)
	}
	defer client.Stop()

	return exec(client)
}

func (c *CountUsersCmd) Usage() {
	fmt.Fprintf(c.flagSet.Output(), "  %s\tCount users\n", cmdName)
}

func newCountUsersCmd() cmd.Cmd {
	result := &CountUsersCmd{}
	result.flagSet = flag.NewFlagSet(cmdName, flag.ExitOnError)
	return result
}

func init() {
	cmd.Add(cmdName, newCountUsersCmd)
}
