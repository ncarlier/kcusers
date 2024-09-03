package version

import (
	"fmt"
	"runtime/debug"
	"time"
)

// Version of the app
var Version = "snapshot"

// GitCommit is the GIT commit revision
var GitCommit = "n/a"

// Built is the built date
var Built = "n/a"

// Print version to stdout
func Print() {
	fmt.Printf(`Version:    %s
Git commit: %s
Built:      %s

Copyright (C) 2024  Nicolas Carlier
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
`, Version, GitCommit, Built)
}

func init() {
	if Version != "snapshot" {
		return
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	Version = info.Main.Version
	for _, kv := range info.Settings {
		if kv.Value == "" {
			continue
		}
		switch kv.Key {
		case "vcs.revision":
			GitCommit = kv.Value[:7]
		case "vcs.time":
			lastCommit, _ := time.Parse(time.RFC3339, kv.Value)
			Built = lastCommit.Format(time.RFC1123)
		}
	}
}
