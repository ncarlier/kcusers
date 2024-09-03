package countsessions

import (
	"errors"
	"flag"
	"fmt"

	"github.com/ncarlier/kcusers/cmd"
	"github.com/ncarlier/kcusers/internal/config"
	"github.com/ncarlier/kcusers/pkg/keycloak"
	"github.com/ncarlier/kcusers/pkg/uuid"
)

const cmdName = "count-sessions"

type CountSessionsCmd struct {
	flagSet *flag.FlagSet
	cid     string
}

func (c *CountSessionsCmd) Exec(args []string, conf *config.Config) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}

	client, err := keycloak.NewKeycloakClient(&conf.Keycloak)
	if err != nil {
		return fmt.Errorf("unable to create Keycloak client: %w", err)
	}
	defer client.Stop()

	if c.cid == "" {
		return errors.New("client ID is required")
	}

	if !uuid.IsUUID(c.cid) {
		return errors.New("invalid client ID")
	}

	return exec(client, c.cid)
}

func (c *CountSessionsCmd) Usage() {
	cmd.PrintCmdUsage(c.flagSet.Output(), cmdName, "Count sessions")
}

func newCountSessionsCmd() cmd.Cmd {
	result := &CountSessionsCmd{}
	result.flagSet = flag.NewFlagSet(cmdName, flag.ExitOnError)
	result.flagSet.StringVar(&result.cid, "cid", "", "Client ID (not client-id)")
	return result
}

func init() {
	cmd.Add(cmdName, newCountSessionsCmd)
}
