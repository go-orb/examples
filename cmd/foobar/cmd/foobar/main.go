// Package main implements the foobar example app
package main

import (
	"fmt"
	"os"

	"github.com/go-orb/examples/cmd/foobar/cmd/foobar/pkg/monolith"
	service1 "github.com/go-orb/examples/cmd/foobar/service/service1/pkg/service"
	"github.com/go-orb/go-orb/cli"

	_ "github.com/go-orb/plugins/client/middleware/log"
	_ "github.com/go-orb/plugins/client/orb"
	_ "github.com/go-orb/plugins/client/orb_transport/grpc"
	_ "github.com/go-orb/plugins/codecs/goccyjson"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/slog"
	_ "github.com/go-orb/plugins/registry/memory"
)

func main() {
	app := cli.App{
		Name:     "foobar",
		Version:  "unset",
		Usage:    "A foobar example app",
		NoAction: true,
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

	app.Commands = append(app.Commands, monolith.MainCommands()...)
	app.Commands = append(app.Commands, service1.MainCommands()...)

	appContext := cli.NewAppContext(&app)

	_, err := run(appContext, os.Args)
	if err != nil {
		fmt.Printf("run error: %s\n", err)
		os.Exit(1)
	}
}
