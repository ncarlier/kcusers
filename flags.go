package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ncarlier/kcusers/cmd"
)

var (
	// ConfigFile is the flag used to load the config file
	configFile string
)

// fisrtCommand restun first command of argument list
func fisrtCommand(args []string) (name string, index int) {
	for idx, arg := range args {
		if strings.HasPrefix(arg, "-") {
			// ignore flags
			continue
		}
		return arg, idx
	}
	return "", -1
}

func init() {
	// -config
	flag.StringVar(&configFile, "config", "", "config file")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s <OPTIONS> <COMMAND>\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\nAvailable options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nAvailable commands:\n")
		for _, c := range cmd.Commands {
			c.Usage()
		}
	}
}
