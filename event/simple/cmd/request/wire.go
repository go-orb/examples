//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"fmt"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/event"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/plugins/cli/urfave"
	"github.com/go-orb/wire"
)

// wireRunResult is here so "wire" has a type for the return value of wireRun.
type wireRunResult struct{}

// wireRunCallback is the actual code that runs the business logic.
type wireRunCallback func(
	ctx context.Context,
	logger log.Logger,
	event event.Handler,
) error

func wireRun(
	serviceContext *cli.ServiceContext,
	components *types.Components,
	logger log.Logger,
	eventHandler event.Handler,
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

	//
	// Actual code
	runErr := cb(serviceContext.Context(), logger, eventHandler)

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

		log.ProvideNoOpts,

		event.ProvideNoOpts,

		wireRun,
	))
}
