// Package main contains a fake login server.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/server"
	"github.com/go-orb/go-orb/types"

	authHandler "github.com/go-orb/examples/rest/auth/handler/auth"
	authProto "github.com/go-orb/examples/rest/auth/proto/auth"

	mdrpc "github.com/go-orb/plugins/server/drpc"

	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/log/slog"

	_ "github.com/go-orb/plugins-experimental/registry/mdns"
)

// provideServerOpts provides options for the go-orb server and registers handlers.
//
//nolint:unparam
func provideServerOpts(logger log.Logger) ([]server.ConfigOption, error) {
	opts := []server.ConfigOption{}

	// Register server Handlers.
	hInstance := authHandler.New([]byte("thisIsAWellKnownSecretItCallsForHackMe"), logger)
	hRegister := authProto.RegisterAuthHandler(hInstance)
	server.Handlers.Add(authProto.HandlerAuth, hRegister)

	opts = append(opts, server.WithEntrypointConfig(mdrpc.NewConfig(
		mdrpc.WithName("drpc"),
		mdrpc.WithHandlers(hRegister),
	)))

	logger.Info("Started")

	return opts, nil
}

func main() {
	var (
		serviceName    = types.ServiceName("orb.examples.rest.auth.server")
		serviceVersion = types.ServiceVersion("v0.0.1")
	)

	components, err := newComponents(serviceName, serviceVersion)
	if err != nil {
		log.Error("while creating components", "err", err)
		os.Exit(1)
	}

	for _, c := range components {
		err := c.Start()
		if err != nil {
			log.Error("Failed to start", err, "component", c.Type())
			os.Exit(1)
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Blocks until we get a sigint/sigterm
	<-done

	// Shutdown.
	ctx := context.Background()

	for k := range components {
		c := components[len(components)-1-k]

		err := c.Stop(ctx)
		if err != nil {
			log.Error("Failed to stop", err, "component", c.Type())
		}
	}
}
