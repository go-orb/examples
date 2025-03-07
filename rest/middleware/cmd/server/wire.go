//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"context"
	"fmt"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/server"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/plugins/cli/urfave"

	"github.com/go-orb/wire"

	echoHandler "github.com/go-orb/examples/rest/middleware/handler/echo"
	echoProto "github.com/go-orb/examples/rest/middleware/proto/echo"
)

// provideServerOpts provides options for the go-orb server and registers handlers.
//
//nolint:unparam
func provideServerOpts(logger log.Logger) ([]server.ConfigOption, error) {
	opts := []server.ConfigOption{}

	// Register server Handlers.
	hInstance := echoHandler.New(logger)
	hRegister := echoProto.RegisterEchoHandler(hInstance)
	server.Handlers.Add(echoProto.HandlerEcho, hRegister)

	return opts, nil
}

// wireRunResult is here so "wire" has a type for the return value of wireRun.
// wire needs a explicit type for each provider including "wireRun".
type wireRunResult struct{}

// wireRunCallback is the actual code that runs the business logic.
type wireRunCallback func(
	ctx context.Context,
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	logger log.Logger,
) error

func wireRun(
	serviceContext *cli.ServiceContext,
	components *types.Components,
	logger log.Logger,
	server server.Server,
) (wireRunResult, error) {
	// Orb start
	for _, c := range components.Iterate(false) {
		logger.Debug("Starting", "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))

		err := c.Start(serviceContext.Context())
		if err != nil {
			logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			return wireRunResult{}, fmt.Errorf("failed to start component %s/%s: %w", c.Type(), c.String(), err)
		}
	}

	// Blocks until interrupt
	<-serviceContext.Context().Done()

	// Orb shutdown.
	ctx := context.Background()

	for _, c := range components.Iterate(true) {
		logger.Debug("Stopping", "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))

		err := c.Stop(ctx)
		if err != nil {
			logger.Error("Failed to stop", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
		}
	}

	return wireRunResult{}, nil
}

func run(
	appContext *cli.AppContext,
	args []string,
) (wireRunResult, error) {
	panic(wire.Build(
		urfave.ProvideParser,
		cli.ProvideParsedFlagsFromArgs,

		cli.ProvideSingleServiceContext,
		types.ProvideComponents,

		cli.ProvideConfigData,
		cli.ProvideServiceName,
		cli.ProvideServiceVersion,

		log.ProvideNoOpts,
		registry.ProvideNoOpts,

		provideServerOpts,
		server.Provide,

		wireRun,
	))
}
