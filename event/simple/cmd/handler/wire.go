//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"fmt"

	"github.com/go-orb/examples/event/simple/pb/user_new"
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/event"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/plugins/cli/urfave"
	"github.com/go-orb/wire"
	"github.com/lithammer/shortuuid/v3"
)

// wireRunResult is here so "wire" has a type for the return value of wireRun.
type wireRunResult struct{}

func wireRun(
	serviceContext *cli.ServiceContext,
	components *types.Components,
	logger log.Logger,
	eventHandler event.Type,
) (wireRunResult, error) {
	// Orb start
	for _, c := range components.Iterate(false) {
		logger.Debug("Starting", "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))

		err := c.Start(serviceContext.Context())
		if err != nil {
			logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			return wireRunResult{}, fmt.Errorf("failed to start component %s/%s: %w", c.Type(), c.String(), err)
		}
	}

	//
	// Actual code
	userNewHandler := func(_ context.Context, req *user_new.Request) (*user_new.Resp, error) {
		return &user_new.Resp{Name: req.GetName(), Uuid: shortuuid.New()}, nil
	}

	event.HandleRequest(serviceContext.Context(), eventHandler, "user.new", userNewHandler)
	// Blocks until interrupt
	<-serviceContext.Context().Done()

	// Orb shutdown.
	ctx := context.Background()

	for _, c := range components.Iterate(true) {
		logger.Debug("Stopping", "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))

		err := c.Stop(ctx)
		if err != nil {
			logger.Error("Failed to stop", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
		}
	}

	return wireRunResult{}, nil
}

func run(
	appContext *cli.AppContext,
	args []string,
) (wireRunResult, error) {
	panic(wire.Build(
		urfave.ProvideParser,
		cli.ProvideParsedFlagsFromArgs,

		cli.ProvideSingleServiceContext,
		types.ProvideComponents,

		cli.ProvideAppConfigData,
		cli.ProvideServiceConfigData,

		log.ProvideNoOpts,

		event.ProvideNoOpts,

		wireRun,
	))
}
