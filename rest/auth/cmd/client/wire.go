//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"fmt"

	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/wire"
)

// wireRunResult is here so "wire" has a type for the return value of wireRun.
// wire needs a explicit type for each provider including "wireRun".
type wireRunResult string

// wireRunCallback is the actual code that runs the business logic.
type wireRunCallback func(
	ctx context.Context,
	logger log.Logger,
	client client.Type,
) error

func wireRun(
	ctx context.Context,
	logger log.Logger,
	client client.Type,
	components *types.Components,
	cb wireRunCallback,
) (wireRunResult, error) {
	// Orb start
	for _, c := range components.Iterate(false) {
		err := c.Start(ctx)
		if err != nil {
			logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			return "", err
		}
	}

	//
	// Actual code
	runErr := cb(ctx, logger, client)

	// Orb shutdown.
	ctx = context.Background()

	for _, c := range components.Iterate(true) {
		err := c.Stop(ctx)
		if err != nil {
			logger.Error("Failed to stop", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
		}
	}

	return "", runErr
}

// run combines everything above and
func run(
	ctx context.Context,
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	cb wireRunCallback,
) (wireRunResult, error) {
	panic(wire.Build(
		types.ProvideComponents,
		wire.Value(types.ConfigData{}),
		provideLoggerOpts,
		log.Provide,
		registry.ProvideNoOpts,
		provideClientOpts,
		client.Provide,
		wireRun,
	))
}
