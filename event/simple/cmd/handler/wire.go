//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-orb/go-orb/event"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/plugins/config/source/cli/urfave"
	"github.com/go-orb/wire"
)

type wireRunResult string

type wireRunCallback func(
	event event.Handler,
	done chan os.Signal,
) error

func wireRun(
	logger log.Logger,
	event event.Handler,
	components *types.Components,
	cb wireRunCallback,
) (wireRunResult, error) {
	// Orb start
	for _, c := range components.Iterate(false) {
		err := c.Start()
		if err != nil {
			logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			return "", err
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	//
	// Actual code
	runErr := cb(event, done)

	// Orb shutdown.
	ctx := context.Background()

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
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	cb wireRunCallback,
) (wireRunResult, error) {
	panic(wire.Build(
		types.ProvideComponents,
		urfave.ProvideConfigData,
		wire.Value([]log.Option{}),
		log.Provide,
		wire.Value([]event.Option{}),
		event.Provide,
		wireRun,
	))
}
