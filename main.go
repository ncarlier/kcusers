package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ncarlier/kcusers/cmd"
	_ "github.com/ncarlier/kcusers/cmd/all"
	"github.com/ncarlier/kcusers/internal"
	"github.com/ncarlier/kcusers/internal/config"
	"github.com/ncarlier/kcusers/pkg/logger"
)

func main() {
	flag.Parse()

	// show version if asked
	if *internal.ShowVersionFlag {
		internal.PrintVersion()
		os.Exit(0)
	}

	conf := config.NewConfig()
	if configFile != "" {
		if err := conf.LoadConfig(configFile); err != nil {
			slog.Error("unable to load configuration file", "error", err)
			os.Exit(1)
		}
	}
	logger.Configure(conf.Log.Level, conf.Log.Format)

	args := flag.Args()
	commandLabel, idx := fisrtCommand(args)

	if command, ok := cmd.Commands[commandLabel]; ok {
		if err := command.Exec(args[idx+1:], conf); err != nil {
			slog.Error("error during command execution", "error", err, "command", commandLabel)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "undefined command: %s\n", commandLabel)
		flag.Usage()
		os.Exit(0)
	}
}
