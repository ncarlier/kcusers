package deleteuser

import (
	"errors"
	"flag"
	"fmt"

	"github.com/ncarlier/kcusers/cmd"
	"github.com/ncarlier/kcusers/internal/config"
	"github.com/ncarlier/kcusers/pkg/keycloak"
)

const cmdName = "delete-users"

type execParams struct {
	filename  string
	concurent uint
	dryRun    bool
}

type DeleteUsersCmd struct {
	execParams
	flagSet *flag.FlagSet
}

func (c *DeleteUsersCmd) Exec(args []string, conf *config.Config) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}

	client, err := keycloak.NewKeycloakClient(&conf.Keycloak)
	if err != nil {
		return fmt.Errorf("unable to create Keycloak client: %w", err)
	}
	defer client.Stop()

	if c.filename == "" {
		return errors.New("filename is required")
	}
	if c.concurent == 0 || c.concurent > 100 {
		return errors.New("invalid concurent value")
	}

	return exec(client, c.execParams)
}

func (c *DeleteUsersCmd) Usage() {
	cmd.PrintCmdUsage(c.flagSet.Output(), cmdName, "Delete users using UID list")
}

func newDeleteUsersCmd() cmd.Cmd {
	result := &DeleteUsersCmd{}
	result.flagSet = flag.NewFlagSet(cmdName, flag.ExitOnError)
	result.flagSet.StringVar(&result.filename, "f", "", "List of users to delete")
	result.flagSet.UintVar(&result.concurent, "concurent", 5, "Concurent API calls (0..100)")
	result.flagSet.BoolVar(&result.dryRun, "dry-run", true, "Dry run execution")
	return result
}

func init() {
	cmd.Add(cmdName, newDeleteUsersCmd)
}
