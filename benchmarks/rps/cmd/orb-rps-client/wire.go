//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/plugins/cli/urfave"

	"github.com/go-orb/wire"
)

type wireRunCallback func(
	ctx context.Context,
	cfg *clientConfig,
	logger log.Logger,
	cli client.Type,
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
	// Orb start
	for _, c := range components.Iterate(false) {
		logger.Debug("Starting", "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))

		err := c.Start(serviceContext.Context())
		if err != nil {
			logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			return wireRunResult{}, fmt.Errorf("failed to start component %s/%s: %w", c.Type(), c.String(), err)
		}
	}

	// Actual code
	runErr := cb(serviceContext.Context(), cfg, logger, clientWire)

	// Orb shutdown.
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

	if err := config.Parse([]string{configSection}, configs, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func provideClientOpts(cfg *clientConfig) ([]client.Option, error) {
	return []client.Option{client.WithClientPoolHosts(1), client.WithClientPoolSize(cfg.PoolSize), client.WithClientConnectionTimeout(time.Duration(cfg.Timeout) * time.Second)}, nil
}

func run(
	appContext *cli.AppContext,
	args []string,
	cb wireRunCallback,
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

		provideClientConfig,
		provideClientOpts,
		client.Provide,

		wireRun,
	))
}
