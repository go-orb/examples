// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/go-orb/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"fmt"
	"github.com/go-orb/examples/benchmarks/rps/handler/echo"
	echo2 "github.com/go-orb/examples/benchmarks/rps/proto/echo"
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/server"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/plugins/cli/urfave"
	"github.com/go-orb/plugins/server/drpc"
	"github.com/go-orb/plugins/server/grpc"
	"github.com/go-orb/plugins/server/http"
)

import (
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

// Injectors from wire.go:

func run(appContext *cli.AppContext, args []string) (wireRunResult, error) {
	serviceContext, err := cli.ProvideSingleServiceContext(appContext)
	if err != nil {
		return wireRunResult{}, err
	}
	v, err := types.ProvideComponents()
	if err != nil {
		return wireRunResult{}, err
	}
	serviceName, err := cli.ProvideServiceName(serviceContext)
	if err != nil {
		return wireRunResult{}, err
	}
	parserFunc, err := urfave.ProvideParser()
	if err != nil {
		return wireRunResult{}, err
	}
	v2, err := cli.ProvideParsedFlagsFromArgs(appContext, parserFunc, args)
	if err != nil {
		return wireRunResult{}, err
	}
	configData, err := cli.ProvideConfigData(serviceContext, v2)
	if err != nil {
		return wireRunResult{}, err
	}
	logger, err := log.ProvideNoOpts(serviceName, configData, v)
	if err != nil {
		return wireRunResult{}, err
	}
	serviceVersion, err := cli.ProvideServiceVersion(serviceContext)
	if err != nil {
		return wireRunResult{}, err
	}
	registryType, err := registry.ProvideNoOpts(serviceName, serviceVersion, configData, v, logger)
	if err != nil {
		return wireRunResult{}, err
	}
	v3, err := provideServerOpts()
	if err != nil {
		return wireRunResult{}, err
	}
	serverServer, err := server.Provide(serviceName, configData, v, logger, registryType, v3...)
	if err != nil {
		return wireRunResult{}, err
	}
	mainWireRunResult, err := wireRun(serviceContext, v, logger, serverServer)
	if err != nil {
		return wireRunResult{}, err
	}
	return mainWireRunResult, nil
}

// wire.go:

// provideServerOpts provides options for the go-orb server.
func provideServerOpts() ([]server.ConfigOption, error) {

	hInstance := new(echo.Handler)
	hRegister := echo2.RegisterEchoHandler(hInstance)

	opts := []server.ConfigOption{}
	opts = append(opts, server.WithEntrypointConfig(grpc.NewConfig(grpc.WithName("grpc"), grpc.WithInsecure(), grpc.WithHandlers(hRegister))))
	opts = append(opts, server.WithEntrypointConfig(http.NewConfig(http.WithName("http"), http.WithInsecure(), http.WithHandlers(hRegister))))
	opts = append(opts, server.WithEntrypointConfig(http.NewConfig(http.WithName("https"), http.WithDisableHTTP2(), http.WithHandlers(hRegister))))
	opts = append(opts, server.WithEntrypointConfig(http.NewConfig(http.WithName("h2c"), http.WithInsecure(), http.WithAllowH2C(), http.WithHandlers(hRegister))))
	opts = append(opts, server.WithEntrypointConfig(http.NewConfig(http.WithName("http2"), http.WithHandlers(hRegister))))
	opts = append(opts, server.WithEntrypointConfig(http.NewConfig(http.WithName("http3"), http.WithHTTP3(), http.WithHandlers(hRegister))))
	opts = append(opts, server.WithEntrypointConfig(drpc.NewConfig(drpc.WithName("drpc"), drpc.WithHandlers(hRegister))))

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

	for _, c := range components.Iterate(false) {
		logger.Debug("Starting", "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))

		err := c.Start(serviceContext.Context())
		if err != nil {
			logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			return wireRunResult{}, fmt.Errorf("failed to start component %s/%s: %w", c.Type(), c.String(), err)
		}
	}

	<-serviceContext.Context().Done()

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
