//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/plugins/cli/urfave"
	"github.com/go-orb/wire"

	"github.com/go-orb/examples/realworld/cmd/realworld/pkg/monolith"
)

type wireRunResult struct{}

func wireRun(
	appContext *cli.AppContext,
	appConfigData cli.AppConfigData,
	monolithRunner monolith.Runner,
) (wireRunResult, error) {
	return wireRunResult{}, monolithRunner()
}

func run(
	appContext *cli.AppContext,
	args []string,
) (wireRunResult, error) {
	panic(wire.Build(
		urfave.ProvideParser,
		cli.ProvideParsedFlagsFromArgs,

		cli.ProvideAppConfigData,

		monolith.ProvideRunner,

		wireRun,
	))
}
