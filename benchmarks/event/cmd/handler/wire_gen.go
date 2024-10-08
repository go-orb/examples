// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/event"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

import (
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/config/source/cli/urfave"
	_ "github.com/go-orb/plugins/event/natsjs"
	_ "github.com/go-orb/plugins/log/slog"
)

// Injectors from wire.go:

// run combines everything above and
func run(serviceName types.ServiceName, serviceVersion types.ServiceVersion, cb wireRunCallback) (wireRunResult, error) {
	configData, err := provideConfigData(serviceName, serviceVersion)
	if err != nil {
		return "", err
	}
	v := _wireValue
	logger, err := log.Provide(serviceName, configData, v...)
	if err != nil {
		return "", err
	}
	v2 := _wireValue2
	eventType, err := event.Provide(serviceName, configData, logger, v2...)
	if err != nil {
		return "", err
	}
	v3, err := provideComponents(logger, eventType)
	if err != nil {
		return "", err
	}
	mainWireRunResult, err := wireRun(serviceName, v3, configData, logger, eventType, cb)
	if err != nil {
		return "", err
	}
	return mainWireRunResult, nil
}

var (
	_wireValue  = []log.Option{}
	_wireValue2 = []event.Option{}
)

// wire.go:

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

// provideComponents creates a slice of components out of the arguments.
func provideComponents(
	logger log.Logger, event2 event.Type,

) ([]types.Component, error) {
	components := []types.Component{}
	components = append(components, logger)
	components = append(components, event2)

	return components, nil
}

type wireRunResult string

type wireRunCallback func(
	serviceName types.ServiceName,
	configs types.ConfigData,
	logger log.Logger, event2 event.Type,

	done chan os.Signal,
) error

func wireRun(
	serviceName types.ServiceName,
	components []types.Component,
	configs types.ConfigData,
	logger log.Logger, event2 event.Type,

	cb wireRunCallback,
) (wireRunResult, error) {

	for _, c := range components {
		err := c.Start()
		if err != nil {
			log.Error("Failed to start", err, "component", c.Type())
			os.Exit(1)
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	runErr := cb(serviceName, configs, logger, event2, done)

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
