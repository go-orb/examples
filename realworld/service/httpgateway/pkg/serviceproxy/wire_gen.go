// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/go-orb/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package serviceproxy

import (
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/client"
	service2 "github.com/go-orb/service/httpgateway"
	"github.com/go-orb/service/httpgateway/pkg/service"
)

// Injectors from wire.go:

// ProvideRunner provides the httpgateway runner.
func ProvideRunner(appContext *cli.AppContext, appConfigData cli.AppConfigData, flags []*cli.Flag) (service.Runner, error) {
	serviceContext, err := provideServiceContext(appContext)
	if err != nil {
		return nil, err
	}
	runner, err := service.ProvideRunner(serviceContext, appConfigData, flags)
	if err != nil {
		return nil, err
	}
	return runner, nil
}

// wire.go:

// Name is the service name.
//
//nolint:gochecknoglobals
var Name = "realworld.service.httpgateway"

func provideServiceContext(appContext *cli.AppContext) (*cli.ServiceContext, error) {
	return cli.NewServiceContext(appContext, Name, ""), nil
}

// ProvideClient provides the httpgateway client.
func ProvideClient(clientWire client.Type) (*service2.Client, error) {
	return service2.New(Name, clientWire), nil
}
