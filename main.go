package main

import (
	"os"

	"github.com/nahria/timetrace/cli"
	"github.com/nahria/timetrace/config"
	"github.com/nahria/timetrace/core"
	"github.com/nahria/timetrace/fs"
	"github.com/nahria/timetrace/out"
)

var version = "UNDEFINED"

func main() {
	c, err := config.FromFile()
	if err != nil {
		out.Warn("%s", err.Error())
	}

	filesystem := fs.New(c)
	timetrace := core.New(c, filesystem)

	if err := cli.RootCommand(timetrace, version).Execute(); err != nil {
		out.Err("%s", err.Error())
		os.Exit(1)
	}
}
