//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

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
	svcCtx *cli.ServiceContextWithConfig,
	logger log.Logger,
	eventHandler event.Type,
) error

func wireRun(
	serviceContext *cli.ServiceContextWithConfig,
	components *types.Components,
	logger log.Logger,
	event event.Type,
	cb wireRunCallback,
) (wireRunResult, error) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Orb start
	for _, c := range components.Iterate(false) {
		err := c.Start(ctx)
		if err != nil {
			logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			return wireRunResult{}, err
		}
	}

	// Actual code
	runErr := cb(serviceContext, logger, event)

	// Orb shutdown.
	ctx = context.Background()

	for _, c := range components.Iterate(true) {
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

		cli.ProvideAppConfigData,
		cli.ProvideServiceConfigData,

		log.ProvideNoOpts,

		event.ProvideNoOpts,

		wireRun,
	))
}
