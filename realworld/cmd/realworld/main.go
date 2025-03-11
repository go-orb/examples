// Package main implements the foobar example app
package main

import (
	"fmt"
	"os"

	"github.com/go-orb/examples/realworld/cmd/realworld/pkg/monolith"
	"github.com/go-orb/go-orb/cli"

	httpgatewayserviceproxy "github.com/go-orb/examples/realworld/service/httpgateway/pkg/serviceproxy"
	lobbyservice "github.com/go-orb/examples/realworld/service/lobby/pkg/service"

	_ "github.com/go-orb/plugins-experimental/registry/mdns"
	_ "github.com/go-orb/plugins/client/middleware/log"
	_ "github.com/go-orb/plugins/client/middleware/retry"
	_ "github.com/go-orb/plugins/client/orb"
	_ "github.com/go-orb/plugins/client/orb_transport/grpc"
	_ "github.com/go-orb/plugins/codecs/goccyjson"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/slog"
	_ "github.com/go-orb/plugins/registry/memory"
	_ "github.com/go-orb/plugins/server/grpc"
)

func main() {
	app := cli.App{
		Name:     "realworld",
		Version:  "",
		Usage:    "A realworld example app",
		NoAction: true,
		Flags: []*cli.Flag{
			{
				Name:        "log_level",
				Default:     "INFO",
				EnvVars:     []string{"LOG_LEVEL"},
				ConfigPaths: []cli.FlagConfigPath{{Path: []string{"logger", "level"}}},
				Usage:       "Set the log level, one of TRACE, DEBUG, INFO, WARN, ERROR",
			},
			{
				Name:        "registry",
				Default:     "mdns",
				EnvVars:     []string{"REGISTRY"},
				ConfigPaths: []cli.FlagConfigPath{{Path: []string{"registry", "plugin"}}},
				Usage:       "Set the registry plugin, one of mdns, consul, memory",
			},
		},
		Commands:      []*cli.Command{},
		Configs:       []string{config},
		ConfigsFormat: []string{"yaml"},
	}

	app.Commands = append(app.Commands, monolith.MainCommands()...)
	app.Commands = append(app.Commands, httpgatewayserviceproxy.MainCommands()...)
	app.Commands = append(app.Commands, lobbyservice.MainCommands()...)

	appContext := cli.NewAppContext(&app)

	_, err := run(appContext, os.Args)
	if err != nil {
		//nolint:forbidigo
		fmt.Printf("run error: %s\n", err)
		os.Exit(1)
	}
}
