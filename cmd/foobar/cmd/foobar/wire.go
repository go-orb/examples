//go:build wireinject
// +build wireinject

package main

import (
	"fmt"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/plugins/cli/urfave"
	"github.com/go-orb/wire"

	"github.com/go-orb/examples/cmd/foobar/cmd/foobar/pkg/monolith"
	service1 "github.com/go-orb/examples/cmd/foobar/service/service1/pkg/service"
)

type wireRunResult struct{}

func wireRun(
	appContext *cli.AppContext,
	service1Runner service1.Runner,
	monolithRunner monolith.Runner,
) (wireRunResult, error) {
	var runErr error
	switch appContext.SelectedService {
	case service1.Name:
		runErr = service1Runner(appContext.SelectedCommand)
	case monolith.Name:
		runErr = monolithRunner(appContext.SelectedCommand)
	default:
		runErr = fmt.Errorf("unknown service: %s", appContext.SelectedService)
	}

	return wireRunResult{}, runErr
}

func run(
	appContext *cli.AppContext,
	args []string,
) (wireRunResult, error) {
	panic(wire.Build(
		urfave.ProvideParser,
		cli.ProvideParsedFlagsFromArgs,

		service1.ProvideRunner,
		monolith.ProvideRunner,

		wireRun,
	))
}
