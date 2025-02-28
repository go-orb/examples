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

// provideServerOpts provides options for the go-orb server and registers handlers.
//
//nolint:unparam
func provideServerOpts(logger log.Logger) ([]server.ConfigOption, error) {
	opts := []server.ConfigOption{}

	// Register server Handlers.
	hInstance := authHandler.New([]byte("thisIsAWellKnownSecretItCallsForHackMe"), logger)
	hRegister := authV1Proto.RegisterAuthHandler(hInstance)
	server.Handlers.Add(authV1Proto.HandlerAuth, hRegister)

	opts = append(opts, server.WithEntrypointConfig(mdrpc.NewConfig(
		mdrpc.WithName("drpc"),
		mdrpc.WithHandlers(hRegister),
	)))

	logger.Info("Started")

	return opts, nil
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
