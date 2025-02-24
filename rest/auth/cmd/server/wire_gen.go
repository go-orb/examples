// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/go-orb/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/server"
	"github.com/go-orb/go-orb/types"
)

import (
	_ "github.com/go-orb/plugins-experimental/registry/mdns"
	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/log/slog"
)

// Injectors from wire.go:

// newComponents combines everything above and returns a slice of components.
func newComponents(serviceName types.ServiceName, serviceVersion types.ServiceVersion) ([]types.Component, error) {
	configData := _wireConfigDataValue
	v, err := provideLoggerOpts()
	if err != nil {
		return nil, err
	}
	logger, err := log.Provide(serviceName, configData, v...)
	if err != nil {
		return nil, err
	}
	v2 := _wireValue
	registryType, err := registry.Provide(serviceName, serviceVersion, configData, logger, v2...)
	if err != nil {
		return nil, err
	}
	v3, err := provideServerOpts(logger)
	if err != nil {
		return nil, err
	}
	serverServer, err := server.Provide(serviceName, configData, logger, registryType, v3...)
	if err != nil {
		return nil, err
	}
	v4, err := provideComponents(logger, registryType, serverServer)
	if err != nil {
		return nil, err
	}
	return v4, nil
}

var (
	_wireConfigDataValue = types.ConfigData{}
	_wireValue           = []registry.Option{}
)

// wire.go:

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
