//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/server"
	"github.com/go-orb/go-orb/types"

	"github.com/go-orb/wire"
)

// provideLoggerOpts returns the logger options.
func provideLoggerOpts() ([]log.Option, error) {
	return []log.Option{log.WithLevel("TRACE")}, nil
}

// provideComponents creates a slice of components out of the arguments.
func provideComponents(
	logger log.Logger,
	reg registry.Type,
	srv server.Server,
) ([]types.Component, error) {
	components := []types.Component{}
	components = append(components, logger)
	components = append(components, reg)
	components = append(components, &srv)

	return components, nil
}

// newComponents combines everything above and returns a slice of components.
func newComponents(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
) ([]types.Component, error) {
	panic(wire.Build(
		wire.Value(types.ConfigData{}),
		provideLoggerOpts,
		log.Provide,
		wire.Value([]registry.Option{}),
		registry.Provide,
		provideServerOpts,
		server.Provide,
		provideComponents,
	))
}
