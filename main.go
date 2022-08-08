package main

import (
	"flag"
	"fmt"
	loggerpkg "github.com/reneManqueros/logger"
	"log"
	"ratatoskr/cmd"
	"ratatoskr/internal/models/config"
	"runtime/debug"
)

func getBuildInfo() string {
	info, _ := debug.ReadBuildInfo()
	vcsRev := ""
	vcsTime := ""
	vcsDirty := ""
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			vcsRev = s.Value
		}
		if s.Key == "vcs.modified" {
			vcsDirty = "-Dirty"
		}
		if s.Key == "vcs.time" {
			vcsTime = s.Value
		}
	}
	return fmt.Sprintf("Revision: %s%s\nFrom: %s\n", vcsRev, vcsDirty, vcsTime)
}

func init() {
	isDebugParameter := flag.String("debug", "", "-mode=debug")
	flag.Parse()
	err := config.Settings.Load()
	if err != nil {
		log.Fatal(err)
	}
	loggerpkg.LogLevel = "INFO"
	if *isDebugParameter == "debug" {
		loggerpkg.LogLevel = "DEBUG"
	}
}

func main() {
	loggerpkg.Logger{}.Debug(getBuildInfo())
	cmd.Execute()
}
