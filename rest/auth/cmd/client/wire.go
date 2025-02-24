//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"os"

	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/wire"
)

// provideLoggerOpts returns the logger options.
func provideLoggerOpts() ([]log.Option, error) {
	return []log.Option{log.WithLevel("TRACE")}, nil
}

func provideClientOpts() ([]client.Option, error) {
	return []client.Option{client.WithClientMiddleware(client.MiddlewareConfig{Name: "log"})}, nil
}

// provideComponents creates a slice of components out of the arguments.
func provideComponents(
	logger log.Logger,
	client client.Type,
) ([]types.Component, error) {
	components := []types.Component{}
	components = append(components, logger)
	components = append(components, client)

	return components, nil
}

type wireRunResult string

type wireRunCallback func(
	logger log.Logger,
	client client.Type,
) error

func wireRun(
	_ types.ServiceName,
	components []types.Component,
	_ types.ConfigData,
	logger log.Logger,
	client client.Type,
	cb wireRunCallback,
) (wireRunResult, error) {
	//
	// Orb start
	for _, c := range components {
		err := c.Start()
		if err != nil {
			log.Error("Failed to start", err, "component", c.Type())
			os.Exit(1)
		}
	}

	//
	// Actual code
	runErr := cb(logger, client)

	//
	// Orb shutdown.
	ctx := context.Background()

	for k := range components {
		c := components[len(components)-1-k]

		err := c.Stop(ctx)
		if err != nil {
			log.Error("Failed to stop", err, "component", c.Type())
		}
	}

	return "", runErr
}

// run combines everything above and
func run(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	cb wireRunCallback,
) (wireRunResult, error) {
	panic(wire.Build(
		wire.Value([]types.ConfigData{}),
		provideLoggerOpts,
		log.Provide,
		wire.Value([]registry.Option{}),
		registry.Provide,
		provideClientOpts,
		client.Provide,
		provideComponents,
		wireRun,
	))
}
