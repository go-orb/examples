//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/server"
	"github.com/go-orb/go-orb/types"

	"github.com/go-orb/examples/benchmarks/rps/handler/echo"
	proto "github.com/go-orb/examples/benchmarks/rps/proto/echo"
	"github.com/go-orb/plugins/server/drpc"
	mgrpc "github.com/go-orb/plugins/server/grpc"
	mhttp "github.com/go-orb/plugins/server/http"

	"github.com/go-orb/wire"
)

// provideConfigData reads the config from cli and returns it.
func provideConfigData(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
) (types.ConfigData, error) {
	u, err := url.Parse("cli://urfave")
	if err != nil {
		return nil, err
	}

	cfgSections := types.SplitServiceName(serviceName)

	data, err := config.Read([]*url.URL{u}, cfgSections)

	return data, err
}

// provideServerOpts provides options for the go-orb server.
func provideServerOpts() ([]server.ConfigOption, error) {

	hInstance := new(echo.Handler)
	hRegister := proto.RegisterEchoHandler(hInstance)

	opts := []server.ConfigOption{}
	opts = append(opts, server.WithEntrypointConfig(mgrpc.NewConfig(
		mgrpc.WithName("grpc"),
		mgrpc.WithInsecure(),
		mgrpc.WithHandlers(hRegister),
	)))
	opts = append(opts, server.WithEntrypointConfig(mhttp.NewConfig(
		mhttp.WithName("http"),
		mhttp.WithInsecure(),
		mhttp.WithHandlers(hRegister),
	)))
	opts = append(opts, server.WithEntrypointConfig(mhttp.NewConfig(
		mhttp.WithName("https"),
		mhttp.WithDisableHTTP2(),
		mhttp.WithHandlers(hRegister),
	)))
	opts = append(opts, server.WithEntrypointConfig(mhttp.NewConfig(
		mhttp.WithName("h2c"),
		mhttp.WithInsecure(),
		mhttp.WithAllowH2C(),
		mhttp.WithHandlers(hRegister),
	)))
	opts = append(opts, server.WithEntrypointConfig(mhttp.NewConfig(
		mhttp.WithName("http2"),
		mhttp.WithHandlers(hRegister),
	)))
	opts = append(opts, server.WithEntrypointConfig(mhttp.NewConfig(
		mhttp.WithName("http3"),
		mhttp.WithHTTP3(),
		mhttp.WithHandlers(hRegister),
	)))
	opts = append(opts, server.WithEntrypointConfig(drpc.NewConfig(
		drpc.WithName("drpc"),
		drpc.WithHandlers(hRegister),
	)))

	return opts, nil
}

// wireRunResult is here so "wire" has a type for the return value of wireRun.
// wire needs a explicit type for each provider including "wireRun".
type wireRunResult string

// wireRunCallback is the actual code that runs the business logic.
type wireRunCallback func(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	logger log.Logger,
	done chan os.Signal,
) error

func wireRun(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	logger log.Logger,
	_ server.Server,
	cb wireRunCallback,
) (wireRunResult, error) {
	// Orb start
	for _, c := range types.Components.Iterate(false) {
		err := c.Start()
		if err != nil {
			logger.Error("Failed to start", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			os.Exit(1)
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	//
	// Actual code
	runErr := cb(serviceName, serviceVersion, logger, done)

	// Orb shutdown.
	ctx := context.Background()

	for _, c := range types.Components.Iterate(true) {
		err := c.Stop(ctx)
		if err != nil {
			logger.Error("Failed to stop", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
		}
	}

	return "", runErr
}

// run combines everything above and runs the callback.
func run(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	cb wireRunCallback,
) (wireRunResult, error) {
	panic(wire.Build(
		provideConfigData,
		wire.Value([]log.Option{}),
		log.Provide,
		wire.Value([]registry.Option{}),
		registry.Provide,
		provideServerOpts,
		server.Provide,
		wireRun,
	))
}
