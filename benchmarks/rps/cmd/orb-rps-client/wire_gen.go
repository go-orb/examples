// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/go-orb/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"fmt"
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/plugins/cli/urfave"
	"time"
)

import (
	_ "github.com/go-orb/plugins-experimental/registry/mdns"
	_ "github.com/go-orb/plugins/client/orb"
	_ "github.com/go-orb/plugins/client/orb_transport/drpc"
	_ "github.com/go-orb/plugins/client/orb_transport/grpc"
	_ "github.com/go-orb/plugins/client/orb_transport/h2c"
	_ "github.com/go-orb/plugins/client/orb_transport/http"
	_ "github.com/go-orb/plugins/client/orb_transport/http3"
	_ "github.com/go-orb/plugins/client/orb_transport/https"
	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/slog"
	_ "github.com/go-orb/plugins/registry/consul"
)

// Injectors from wire.go:

func run(appContext *cli.AppContext, args []string, cb wireRunCallback) (wireRunResult, error) {
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
	mainClientConfig, err := provideClientConfig(serviceName, configData)
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
	v3, err := provideClientOpts(mainClientConfig)
	if err != nil {
		return wireRunResult{}, err
	}
	clientType, err := client.Provide(serviceName, configData, v, logger, registryType, v3...)
	if err != nil {
		return wireRunResult{}, err
	}
	mainWireRunResult, err := wireRun(serviceContext, v, mainClientConfig, logger, clientType, cb)
	if err != nil {
		return wireRunResult{}, err
	}
	return mainWireRunResult, nil
}

// wire.go:

type wireRunCallback func(
	ctx context.Context,
	cfg *clientConfig,
	logger log.Logger, cli2 client.Type,

) error

// wireRunResult is here so "wire" has a type for the return value of wireRun.
type wireRunResult struct{}

func wireRun(
	serviceContext *cli.ServiceContext,
	components *types.Components,
	cfg *clientConfig,
	logger log.Logger,
	clientWire client.Type,
	cb wireRunCallback,
) (wireRunResult, error) {

	for _, c := range components.Iterate(false) {
		logger.Debug("Starting", "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))

		err := c.Start(serviceContext.Context())
		if err != nil {
			logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			return wireRunResult{}, fmt.Errorf("failed to start component %s/%s: %w", c.Type(), c.String(), err)
		}
	}

	runErr := cb(serviceContext.Context(), cfg, logger, clientWire)

	ctx := context.Background()

	for _, c := range components.Iterate(true) {
		logger.Debug("Stopping", "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))

		err := c.Stop(ctx)
		if err != nil {
			logger.Error("Failed to stop", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
		}
	}

	return wireRunResult{}, runErr
}

func provideClientConfig(serviceName types.ServiceName, configs types.ConfigData) (*clientConfig, error) {
	cfg := &clientConfig{
		BypassRegistry: defaultBypassRegistry,
		PoolSize:       defaultPoolSize,
		Connections:    defaultConnections,
		Duration:       defaultDuration,
		Timeout:        defaultTimeout,
		Threads:        defaultThreads,
		Transport:      defaultTransport,
		PackageSize:    defaultPackageSize,
		ContentType:    defaultContentType,
	}
	config.Dump(configs)

	if err := config.Parse([]string{configSection}, configs, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func provideClientOpts(cfg *clientConfig) ([]client.Option, error) {
	return []client.Option{client.WithClientPoolHosts(1), client.WithClientPoolSize(cfg.PoolSize), client.WithClientConnectionTimeout(time.Duration(cfg.Timeout) * time.Second)}, nil
}
