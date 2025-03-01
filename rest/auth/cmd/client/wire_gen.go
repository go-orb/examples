// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/go-orb/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"fmt"
	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/types"
	"os"
	"os/signal"
	"syscall"
)

import (
	_ "github.com/go-orb/plugins-experimental/registry/mdns"
	_ "github.com/go-orb/plugins/client/middleware/log"
	_ "github.com/go-orb/plugins/client/orb"
	_ "github.com/go-orb/plugins/client/orb/transport/drpc"
	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/log/slog"
)

// Injectors from wire.go:

// run combines everything above and
func run(serviceName types.ServiceName, serviceVersion types.ServiceVersion, cb wireRunCallback) (wireRunResult, error) {
	configData := _wireConfigDataValue
	v, err := types.ProvideComponents()
	if err != nil {
		return "", err
	}
	v2, err := provideLoggerOpts()
	if err != nil {
		return "", err
	}
	logger, err := log.Provide(serviceName, configData, v, v2...)
	if err != nil {
		return "", err
	}
	v3 := _wireValue
	registryType, err := registry.Provide(serviceName, serviceVersion, configData, v, logger, v3...)
	if err != nil {
		return "", err
	}
	v4, err := provideClientOpts()
	if err != nil {
		return "", err
	}
	clientType, err := client.Provide(serviceName, configData, v, logger, registryType, v4...)
	if err != nil {
		return "", err
	}
	mainWireRunResult, err := wireRun(logger, clientType, v, cb)
	if err != nil {
		return "", err
	}
	return mainWireRunResult, nil
}

var (
	_wireConfigDataValue = types.ConfigData{}
	_wireValue           = []registry.Option{}
)

// wire.go:

// wireRunResult is here so "wire" has a type for the return value of wireRun.
// wire needs a explicit type for each provider including "wireRun".
type wireRunResult string

// wireRunCallback is the actual code that runs the business logic.
type wireRunCallback func(
	logger log.Logger, client2 client.Type,

) error

func wireRun(
	logger log.Logger, client2 client.Type,

	components *types.Components,
	cb wireRunCallback,
) (wireRunResult, error) {

	for _, c := range components.Iterate(false) {
		err := c.Start()
		if err != nil {
			logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			return "", err
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	runErr := cb(logger, client2)

	ctx := context.Background()

	for _, c := range components.Iterate(true) {
		err := c.Stop(ctx)
		if err != nil {
			logger.Error("Failed to stop", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
		}
	}

	return "", runErr
}
