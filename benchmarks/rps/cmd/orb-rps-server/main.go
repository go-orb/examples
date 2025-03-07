// Package main contains a simple handler/server example for a event run.
package main

import (
	"fmt"
	"os"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/registry"

	_ "github.com/go-orb/plugins-experimental/registry/mdns"
	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/lumberjack"
	_ "github.com/go-orb/plugins/log/slog"
	_ "github.com/go-orb/plugins/registry/consul"
	_ "github.com/go-orb/plugins/server/http/router/chi"
)

func main() {
	app := cli.App{
		Name:     "benchmarks.rps.server",
		Version:  "",
		Usage:    "A rps benchmarking server",
		NoAction: false,
		Flags: []*cli.Flag{
			{
				Name:        "registry",
				Default:     registry.DefaultRegistry,
				EnvVars:     []string{"REGISTRY"},
				ConfigPaths: []cli.FlagConfigPath{{Path: []string{"registry", "plugin"}}},
				Usage:       "Set the registry plugin, one of mdns, consul, memory",
			},
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
		fmt.Printf("run error: %s\n", err)
		os.Exit(1)
	}
}
