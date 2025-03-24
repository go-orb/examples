// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/go-orb/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-orb/examples/realworld/cmd/realworld/pkg/monolith"
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/plugins/cli/urfave"
)

import (
	_ "github.com/go-orb/plugins/client/middleware/log"
	_ "github.com/go-orb/plugins/client/middleware/retry"
	_ "github.com/go-orb/plugins/client/orb"
	_ "github.com/go-orb/plugins/client/orb_transport/grpc"
	_ "github.com/go-orb/plugins/codecs/goccyjson"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/slog"
	_ "github.com/go-orb/plugins/registry/mdns"
	_ "github.com/go-orb/plugins/registry/memory"
	_ "github.com/go-orb/plugins/server/grpc"
)

// Injectors from wire.go:

func run(appContext *cli.AppContext, args []string) (wireRunResult, error) {
	appConfigData, err := cli.ProvideAppConfigData(appContext)
	if err != nil {
		return wireRunResult{}, err
	}
	parserFunc, err := urfave.ProvideParser()
	if err != nil {
		return wireRunResult{}, err
	}
	v, err := cli.ProvideParsedFlagsFromArgs(appContext, parserFunc, args)
	if err != nil {
		return wireRunResult{}, err
	}
	runner, err := monolith.ProvideRunner(appContext, appConfigData, v)
	if err != nil {
		return wireRunResult{}, err
	}
	mainWireRunResult, err := wireRun(appContext, appConfigData, runner)
	if err != nil {
		return wireRunResult{}, err
	}
	return mainWireRunResult, nil
}

// wire.go:

type wireRunResult struct{}

func wireRun(
	appContext *cli.AppContext,
	appConfigData cli.AppConfigData,
	monolithRunner monolith.Runner,
) (wireRunResult, error) {
	return wireRunResult{}, monolithRunner()
}
