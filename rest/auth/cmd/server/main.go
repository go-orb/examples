// Package main contains a fake login server.
package main

import (
	"os"

	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/server"
	"github.com/go-orb/go-orb/types"

	authHandler "github.com/go-orb/examples/rest/auth/handler/auth"
	authV1Proto "github.com/go-orb/examples/rest/auth/proto/auth_v1"

	mdrpc "github.com/go-orb/plugins/server/drpc"

	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/log/slog"

	_ "github.com/go-orb/plugins-experimental/registry/mdns"
)

// provideLoggerOpts returns the logger options.
func provideLoggerOpts() ([]log.Option, error) {
	return []log.Option{log.WithLevel("TRACE")}, nil
}

// provideServerOpts provides options for the go-orb server.
//
//nolint:unparam
func provideServerOpts() ([]server.ConfigOption, error) {
	opts := []server.ConfigOption{}

	opts = append(opts, server.WithEntrypointConfig(mdrpc.NewConfig(
		mdrpc.WithName("drpc"),
	)))

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

func runner(serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	logger log.Logger,
	done chan os.Signal,
) error {
	logger.Info("Started", "name", serviceName, "version", serviceVersion)

	// Blocks until the process receives a signal.
	<-done

	logger.Info("Stopping", "name", serviceName, "version", serviceVersion)

	return nil
}

func main() {
	var (
		serviceName    = types.ServiceName("orb.examples.rest.auth.server")
		serviceVersion = types.ServiceVersion("v0.0.1")
	)

	if _, err := run(serviceName, serviceVersion, runner); err != nil {
		log.Error("while running", "err", err)
	}
}
