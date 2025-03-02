// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/go-orb/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"fmt"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/server"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/plugins/config/source/cli/urfave"
	"os"
	"os/signal"
	"syscall"
)

import (
	_ "github.com/go-orb/plugins-experimental/registry/mdns"
	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/cli/urfave"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/lumberjack"
	_ "github.com/go-orb/plugins/log/slog"
	_ "github.com/go-orb/plugins/registry/consul"
	_ "github.com/go-orb/plugins/server/drpc"
	_ "github.com/go-orb/plugins/server/grpc"
	_ "github.com/go-orb/plugins/server/http"
	_ "github.com/go-orb/plugins/server/http/router/chi"
)

// Injectors from wire.go:

// run combines everything above and runs the callback.
func run(serviceName types.ServiceName, serviceVersion types.ServiceVersion, cb wireRunCallback) (wireRunResult, error) {
	configData, err := urfave.ProvideConfigData(serviceName, serviceVersion)
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
	v3, err := provideServerOpts(logger)
	if err != nil {
		return "", err
	}
	serverServer, err := server.Provide(serviceName, configData, logger, registryType, v3...)
	if err != nil {
		return "", err
	}
	mainWireRunResult, err := wireRun(serviceName, serviceVersion, logger, serverServer, cb)
	if err != nil {
		return "", err
	}
	return mainWireRunResult, nil
}

var (
	_wireValue  = []log.Option{}
	_wireValue2 = []registry.Option{}
)

// wire.go:

// wireRunResult is here so "wire" has a type for the return value of wireRun.
// wire needs a explicit type for each provider including "wireRun".
type wireRunResult string

// wireRunCallback is the actual code that runs the business logic.
type wireRunCallback func(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	logger log.Logger,
	done chan os.Signal,
) error

func wireRun(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	logger log.Logger,
	_ server.Server,
	cb wireRunCallback,
) (wireRunResult, error) {

	for _, c := range types.Components.Iterate(false) {
		err := c.Start()
		if err != nil {
			logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			return "", err
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	runErr := cb(serviceName, serviceVersion, logger, done)

	ctx := context.Background()

	for _, c := range types.Components.Iterate(true) {
		err := c.Stop(ctx)
		if err != nil {
			logger.Error("Failed to stop", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
		}
	}

	return "", runErr
}
