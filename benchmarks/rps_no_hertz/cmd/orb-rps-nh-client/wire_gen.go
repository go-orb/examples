// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

import (
	_ "github.com/go-orb/plugins/client/middleware/log"
	_ "github.com/go-orb/plugins/client/orb"
	_ "github.com/go-orb/plugins/client/orb/transport/drpc"
	_ "github.com/go-orb/plugins/client/orb/transport/grpc"
	_ "github.com/go-orb/plugins/client/orb/transport/h2c"
	_ "github.com/go-orb/plugins/client/orb/transport/http"
	_ "github.com/go-orb/plugins/client/orb/transport/http3"
	_ "github.com/go-orb/plugins/client/orb/transport/https"
	_ "github.com/go-orb/plugins/codecs/jsonpb"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/cli/urfave"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/slog"
	_ "github.com/go-orb/plugins/registry/consul"
	_ "github.com/go-orb/plugins/registry/mdns"
)

// Injectors from wire.go:

// newComponents combines everything above and returns a slice of components.
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
	registryType, err := registry.Provide(serviceName, serviceVersion, configData, logger, v2...)
	if err != nil {
		return "", err
	}
	v3 := _wireValue3
	clientType, err := client.Provide(serviceName, configData, logger, registryType, v3...)
	if err != nil {
		return "", err
	}
	v4, err := provideComponents(serviceName, serviceVersion, configData, logger, registryType, clientType)
	if err != nil {
		return "", err
	}
	mainWireRunResult, err := wireRun(serviceName, v4, configData, logger, clientType, cb)
	if err != nil {
		return "", err
	}
	return mainWireRunResult, nil
}

var (
	_wireValue  = []log.Option{}
	_wireValue2 = []registry.Option{}
	_wireValue3 = []client.Option{}
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
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	cfgData types.ConfigData,
	logger log.Logger,
	reg registry.Type, client2 client.Type,

) ([]types.Component, error) {
	components := []types.Component{}
	components = append(components, logger)
	components = append(components, reg)

	return components, nil
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
	components []types.Component,
	configs types.ConfigData,
	logger log.Logger,
	cli client.Type,
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

	runErr := cb(serviceName, configs, logger, cli)

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
