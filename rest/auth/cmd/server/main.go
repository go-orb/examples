// Package main contains a fake login server.
package main

import (
	"fmt"
	"os"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/server"

	authHandler "github.com/go-orb/examples/rest/auth/handler/auth"
	authV1Proto "github.com/go-orb/examples/rest/auth/proto/auth_v1"

	mdrpc "github.com/go-orb/plugins/server/drpc"

	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/slog"

	_ "github.com/go-orb/plugins/registry/consul"
	_ "github.com/go-orb/plugins/registry/mdns"
)

// provideServerOpts provides options for the go-orb server.
//
//nolint:unparam
func provideServerOpts() ([]server.ConfigOption, error) {
	opts := []server.ConfigOption{}

	opts = append(opts, server.WithEntrypointConfig("drpc", mdrpc.NewConfig()))

	return opts, nil
}

// provideServerConfigured configures the go-orb server(s).
//
//nolint:unparam
func provideServerConfigured(logger log.Logger, srv server.Server) (serverConfigured, error) {
	// Register server Handlers.
	hInstance := authHandler.New([]byte("thisIsAWellKnownSecretItCallsForHackMe"), logger)
	hRegister := authV1Proto.RegisterAuthHandler(hInstance)

	// Add our server handler to all entrypoints.
	srv.GetEntrypoints().Range(func(_ string, entrypoint server.Entrypoint) bool {
		entrypoint.AddHandler(hRegister)

		return true
	})

	return serverConfigured{}, nil
}

func runner(
	svcCtx *cli.ServiceContext,
	logger log.Logger,
) error {
	logger.Info("Started", "name", svcCtx.Name(), "version", svcCtx.Version())

	// Blocks until the process receives a signal.
	<-svcCtx.Context().Done()

	logger.Info("Stopping", "name", svcCtx.Name(), "version", svcCtx.Version())

	return nil
}

func main() {
	app := cli.App{
		Name:     "orb.examples.rest.middleware.server",
		Version:  "",
		Usage:    "An example app",
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

	_, err := run(appContext, os.Args, runner)
	if err != nil {
		//nolint:forbidigo
		fmt.Printf("run error: %s\n", err)
		os.Exit(1)
	}
}
