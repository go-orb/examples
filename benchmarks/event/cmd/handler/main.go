// Package main contains a simple handler/server example for a event run.
package main

import (
	"fmt"
	"os"

	"github.com/go-orb/go-orb/cli"
	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/event/natsjs"
	_ "github.com/go-orb/plugins/log/slog"
)

func main() {
	app := cli.App{
		Name:     "benchmarks.event.handler",
		Version:  "",
		Usage:    "A event benchmarking handler",
		NoAction: false,
		Flags: []*cli.Flag{
			{
				Name:        "log_level",
				Default:     "INFO",
				EnvVars:     []string{"LOG_LEVEL"},
				ConfigPaths: []cli.FlagConfigPath{{Path: []string{"logger", "level"}}},
				Usage:       "Set the log level, one of TRACE, DEBUG, INFO, WARN, ERROR",
			},
		},
		Commands: []*cli.Command{},
	}

	appContext := cli.NewAppContext(&app)

	_, err := run(appContext, os.Args)
	if err != nil {
		//nolint:forbidigo
		fmt.Printf("run error: %s\n", err)
		os.Exit(1)
	}
}
