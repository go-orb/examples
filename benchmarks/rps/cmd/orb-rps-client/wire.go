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

	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/plugins/config/source/cli/urfave"

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

type wireRunResult string

type wireRunCallback func(
	serviceName types.ServiceName,
	configs types.ConfigData,
	logger log.Logger,
	cli client.Type,
) error

func wireRun(
	serviceName types.ServiceName,
	configs types.ConfigData,
	logger log.Logger,
	cli client.Type,
	cb wireRunCallback,
) (wireRunResult, error) {
	// Orb start
	for _, c := range types.Components.Iterate(false) {
		err := c.Start()
		if err != nil {
			logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			return "", err
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	//
	// Actual code
	runErr := cb(serviceName, configs, logger, cli)

	// Orb shutdown.
	ctx := context.Background()

	for _, c := range types.Components.Iterate(true) {
		err := c.Stop(ctx)
		if err != nil {
			logger.Error("Failed to stop", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
		}
	}

	return "", runErr
}

// newComponents combines everything above and returns a slice of components.
func run(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	cb wireRunCallback,
) (wireRunResult, error) {
	panic(wire.Build(
		urfave.ProvideConfigData,
		wire.Value([]log.Option{}),
		log.Provide,
		wire.Value([]registry.Option{}),
		registry.Provide,
		wire.Value([]client.Option{}),
		client.Provide,
		wireRun,
	))
}
