package initconfig

import (
	"flag"

	"github.com/ncarlier/kcusers/cmd"
	"github.com/ncarlier/kcusers/internal/config"
)

const cmdName = "init-config"

type InitConfigCmd struct {
	filename string
	flagSet  *flag.FlagSet
}

func (c *InitConfigCmd) Exec(args []string, conf *config.Config) error {
	if err := c.flagSet.Parse(args); err != nil {
		return err
	}
	return writeDefaultConfigFile(c.filename)
}

func (c *InitConfigCmd) Usage() {
	cmd.PrintCmdUsage(c.flagSet.Output(), cmdName, "Initialize configuration file")
}

func newInitConfigCmd() cmd.Cmd {
	result := &InitConfigCmd{}

	result.flagSet = flag.NewFlagSet(cmdName, flag.ExitOnError)
	result.flagSet.StringVar(&result.filename, "f", "config.toml", "Configuration file to create")

	return result
}

func init() {
	cmd.Add(cmdName, newInitConfigCmd)
}
