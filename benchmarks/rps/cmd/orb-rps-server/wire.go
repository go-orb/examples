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

	_ "github.com/go-orb/plugins/registry/consul"

	"github.com/go-orb/examples/benchmarks/rps/handler/echo"
	proto "github.com/go-orb/examples/benchmarks/rps/proto/echo"
	"github.com/go-orb/plugins/cli/urfave"
	"github.com/go-orb/plugins/server/drpc"
	mgrpc "github.com/go-orb/plugins/server/grpc"
	mhttp "github.com/go-orb/plugins/server/http"

	"github.com/go-orb/wire"
)

// provideServerOpts provides options for the go-orb server.
func provideServerOpts() ([]server.ConfigOption, error) {

	hInstance := new(echo.Handler)
	hRegister := proto.RegisterEchoHandler(hInstance)

	opts := []server.ConfigOption{}
	opts = append(opts, server.WithEntrypointConfig("grpc", mgrpc.NewConfig(
		mgrpc.WithInsecure(),
		mgrpc.WithHandlers(hRegister),
		mgrpc.WithReflection(true),
	)))
	opts = append(opts, server.WithEntrypointConfig("grpcs", mgrpc.NewConfig(
		mgrpc.WithHandlers(hRegister),
		mgrpc.WithReflection(true),
	)))
	opts = append(opts, server.WithEntrypointConfig("http", mhttp.NewConfig(
		mhttp.WithInsecure(),
		mhttp.WithHandlers(hRegister),
	)))
	opts = append(opts, server.WithEntrypointConfig("https", mhttp.NewConfig(
		mhttp.WithDisableHTTP2(),
		mhttp.WithHandlers(hRegister),
	)))
	opts = append(opts, server.WithEntrypointConfig("h2c", mhttp.NewConfig(
		mhttp.WithInsecure(),
		mhttp.WithAllowH2C(),
		mhttp.WithHandlers(hRegister),
	)))
	opts = append(opts, server.WithEntrypointConfig("http2", mhttp.NewConfig(
		mhttp.WithHandlers(hRegister),
	)))
	opts = append(opts, server.WithEntrypointConfig("http3", mhttp.NewConfig(
		mhttp.WithHTTP3(),
		mhttp.WithHandlers(hRegister),
	)))
	opts = append(opts, server.WithEntrypointConfig("drpc", drpc.NewConfig(
		drpc.WithHandlers(hRegister),
	)))

	return opts, nil
}

// wireRunResult is here so "wire" has a type for the return value of wireRun.
type wireRunResult struct{}

func wireRun(
	serviceContext *cli.ServiceContext,
	components *types.Components,
	logger log.Logger,
	_ server.Server,
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

		cli.ProvideAppConfigData,
		cli.ProvideServiceConfigData,

		log.ProvideNoOpts,
		registry.ProvideNoOpts,

		provideServerOpts,
		server.Provide,

		wireRun,
	))
}
