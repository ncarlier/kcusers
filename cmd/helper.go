package cmd

import (
	"fmt"
	"io"
	"strings"
)

// GetFirstCommand return first command of argument list
func GetFirstCommand(args []string) (name string, index int) {
	for idx, arg := range args {
		if strings.HasPrefix(arg, "-") {
			// ignore flags
			continue
		}
		return arg, idx
	}
	return "", -1
}

// PrintCmdUsage print command usage
func PrintCmdUsage(w io.Writer, cmdName, cmdDesc string) {
	fmt.Fprintf(w, "  %s\t%s\n", cmdName, cmdDesc)
}
