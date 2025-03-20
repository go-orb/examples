// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/go-orb/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"fmt"
	"github.com/go-orb/examples/benchmarks/event/pb/echo"
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/event"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/plugins/cli/urfave"
)

import (
	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/event/natsjs"
	_ "github.com/go-orb/plugins/log/slog"
)

// Injectors from wire.go:

func run(appContext *cli.AppContext, args []string) (wireRunResult, error) {
	serviceContext, err := cli.ProvideSingleServiceContext(appContext)
	if err != nil {
		return wireRunResult{}, err
	}
	v, err := types.ProvideComponents()
	if err != nil {
		return wireRunResult{}, err
	}
	appConfigData, err := cli.ProvideAppConfigData(appContext)
	if err != nil {
		return wireRunResult{}, err
	}
	parserFunc, err := urfave.ProvideParser()
	if err != nil {
		return wireRunResult{}, err
	}
	v2, err := cli.ProvideParsedFlagsFromArgs(appContext, parserFunc, args)
	if err != nil {
		return wireRunResult{}, err
	}
	serviceContextHasConfigData, err := cli.ProvideServiceConfigData(serviceContext, appConfigData, v2)
	if err != nil {
		return wireRunResult{}, err
	}
	logger, err := log.ProvideNoOpts(serviceContextHasConfigData, serviceContext, v)
	if err != nil {
		return wireRunResult{}, err
	}
	eventType, err := event.ProvideNoOpts(serviceContext, v, logger)
	if err != nil {
		return wireRunResult{}, err
	}
	mainWireRunResult, err := wireRun(serviceContext, v, logger, eventType)
	if err != nil {
		return wireRunResult{}, err
	}
	return mainWireRunResult, nil
}

// wire.go:

// wireRunResult is here so "wire" has a type for the return value of wireRun.
type wireRunResult struct{}

func wireRun(
	serviceContext *cli.ServiceContext,
	components *types.Components,
	logger log.Logger,
	eventHandler event.Type,
) (wireRunResult, error) {

	for _, c := range components.Iterate(false) {
		logger.Debug("Starting", "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))

		err := c.Start(serviceContext.Context())
		if err != nil {
			logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			return wireRunResult{}, fmt.Errorf("failed to start component %s/%s: %w", c.Type(), c.String(), err)
		}
	}

	echoHandler := func(_ context.Context, req *echo.Req) (*echo.Resp, error) {
		return &echo.Resp{Payload: req.GetPayload()}, nil
	}
	event.HandleRequest(serviceContext.Context(), eventHandler, "echo.echo", echoHandler)

	<-serviceContext.Context().Done()

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
