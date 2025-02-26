// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/go-orb/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"fmt"
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
	_ "github.com/go-orb/plugins-experimental/registry/mdns"
	_ "github.com/go-orb/plugins/client/middleware/log"
	_ "github.com/go-orb/plugins/client/orb"
	_ "github.com/go-orb/plugins/client/orb/transport/drpc"
	_ "github.com/go-orb/plugins/client/orb/transport/grpc"
	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/cli/urfave"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/slog"
	_ "github.com/go-orb/plugins/registry/consul"
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
	registryType, err := registry.Provide(serviceName, serviceVersion, configData, logger, v2...)
	if err != nil {
		return "", err
	}
	v3 := _wireValue3
	clientType, err := client.Provide(serviceName, configData, logger, registryType, v3...)
	if err != nil {
		return "", err
	}
	mainWireRunResult, err := wireRun(logger, clientType, cb)
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

type wireRunResult string

type wireRunCallback func(
	logger log.Logger, client2 client.Type,

) error

func wireRun(
	logger log.Logger, client2 client.Type,

	cb wireRunCallback,
) (wireRunResult, error) {

	for _, c := range types.Components.Iterate(false) {
		err := c.Start()
		if err != nil {
			logger.Error("Failed to start", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			os.Exit(1)
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	runErr := cb(logger, client2)

	ctx := context.Background()

	for _, c := range types.Components.Iterate(true) {
		err := c.Stop(ctx)
		if err != nil {
			logger.Error("Failed to stop", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
		}
	}

	return "", runErr
}
