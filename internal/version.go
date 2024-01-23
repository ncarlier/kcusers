package internal

import (
	"flag"
	"fmt"
)

// Version of the app
var Version = "snapshot"

// GitCommit is the GIT commit revision
var GitCommit = "n/a"

// Built is the built date
var Built = "n/a"

// ShowVersionFlag is the flag used to print version
var ShowVersionFlag = flag.Bool("version", false, "Print version")

// PrintVersion to stdout
func PrintVersion() {
	fmt.Printf(`Version:    %s
Git commit: %s
Built:      %s

readflow  Copyright (C) 2024  Nicolas Carlier

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

There is NO WARRANTY, to the extent permitted by law.
`, Version, GitCommit, Built)
}
