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

Copyright (C) 2024  Nicolas Carlier
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
`, Version, GitCommit, Built)
}
