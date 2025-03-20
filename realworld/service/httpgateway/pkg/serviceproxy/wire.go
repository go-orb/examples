//go:build wireinject
// +build wireinject

// Package serviceproxy proxies the httpgateway service with a custom name.
package serviceproxy

import (
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/wire"

	httpgateway "github.com/go-orb/service/httpgateway"
	httpgatewayservice "github.com/go-orb/service/httpgateway/pkg/service"
)

// Name is the service name.
//
//nolint:gochecknoglobals
var Name = "realworld.service.httpgateway"

func provideServiceContext(appContext *cli.AppContext) (*cli.ServiceContext, error) {
	return cli.NewServiceContext(appContext, Name, ""), nil
}

// ProvideClient provides the httpgateway client.
func ProvideClient(clientWire client.Type) (*httpgateway.Client, error) {
	return httpgateway.New(Name, clientWire), nil
}

// ProvideRunner provides the httpgateway runner.
func ProvideRunner(appContext *cli.AppContext, appConfigData cli.AppConfigData, flags []*cli.Flag) (httpgatewayservice.Runner, error) {
	panic(wire.Build(
		provideServiceContext,
		httpgatewayservice.ProvideRunner,
	))
}
