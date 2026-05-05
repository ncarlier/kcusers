package resetpassword

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/ncarlier/kcusers/cmd"
	"github.com/ncarlier/kcusers/internal/config"
	"github.com/ncarlier/kcusers/pkg/keycloak"
	"github.com/ncarlier/kcusers/pkg/uuid"
)

const cmdName = "reset-password"

type ResetPasswordCmd struct {
	uid       string
	username  string
	password  string
	temporary bool
	flagSet   *flag.FlagSet
}

func (c *ResetPasswordCmd) Exec(args []string, conf *config.Config) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}

	if c.uid == "" && c.username == "" {
		return errors.New("either uid or username must be provided")
	}

	if c.uid != "" && !uuid.IsUUID(c.uid) {
		return errors.New("invalid user ID")
	}

	password := c.password
	if password == "" {
		fmt.Print("Enter new password: ")
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("unable to read password: %w", err)
		}
		fmt.Println()
		password = string(bytePassword)
		if password == "" {
			return errors.New("password cannot be empty")
		}
	}

	client, err := keycloak.NewKeycloakClient(&conf.Keycloak)
	if err != nil {
		return fmt.Errorf("unable to create Keycloak client: %w", err)
	}
	defer client.Stop()

	return exec(client, c.uid, c.username, password, c.temporary)
}

func (c *ResetPasswordCmd) Usage() {
	cmd.PrintCmdUsage(c.flagSet.Output(), cmdName, "Reset user password")
}

func newResetPasswordCmd() cmd.Cmd {
	result := &ResetPasswordCmd{}
	result.flagSet = flag.NewFlagSet(cmdName, flag.ExitOnError)
	result.flagSet.StringVar(&result.uid, "uid", "", "User ID")
	result.flagSet.StringVar(&result.username, "username", "", "Username (if uid is not provided)")
	result.flagSet.StringVar(&result.password, "password", "", "New password (will prompt if empty)")
	result.flagSet.BoolVar(&result.temporary, "temporary", true, "Force user to change password on next login")
	return result
}

func init() {
	cmd.Add(cmdName, newResetPasswordCmd)
}
