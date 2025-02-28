// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/go-orb/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"fmt"
	"github.com/go-orb/go-orb/event"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/plugins/config/source/cli/urfave"
	"os"
	"os/signal"
	"syscall"
)

import (
	_ "github.com/go-orb/plugins/codecs/goccyjson"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/config/source/cli/urfave"
	_ "github.com/go-orb/plugins/event/natsjs"
	_ "github.com/go-orb/plugins/log/slog"
)

// Injectors from wire.go:

// run combines everything above and
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
	handler, err := event.Provide(serviceName, configData, logger, v2...)
	if err != nil {
		return "", err
	}
	mainWireRunResult, err := wireRun(serviceName, configData, logger, handler, cb)
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

// wireRunResult is here so "wire" has a type for the return value of wireRun.
type wireRunResult string

// wireRunCallback is the actual code that runs the business logic.
type wireRunCallback func(
	serviceName types.ServiceName,
	configs types.ConfigData,
	logger log.Logger,
	eventHandler event.Handler,
	done chan os.Signal,
) error

func wireRun(
	serviceName types.ServiceName,
	configs types.ConfigData,
	logger log.Logger, event2 event.Handler,

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

	runErr := cb(serviceName, configs, logger, event2, done)

	ctx := context.Background()

	for _, c := range types.Components.Iterate(true) {
		err := c.Stop(ctx)
		if err != nil {
			logger.Error("Failed to stop", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
		}
	}

	return "", runErr
}
