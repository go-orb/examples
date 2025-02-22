// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-orb/examples/benchmarks/rps/handler/echo"
	echo2 "github.com/go-orb/examples/benchmarks/rps/proto/echo"
	"github.com/go-orb/go-orb/config"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/server"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/plugins/server/drpc"
	"github.com/go-orb/plugins/server/grpc"
	"github.com/go-orb/plugins/server/http"
	"net/url"
)

import (
	_ "github.com/go-orb/plugins-experimental/registry/mdns"
	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/cli/urfave"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/lumberjack"
	_ "github.com/go-orb/plugins/log/slog"
	_ "github.com/go-orb/plugins/registry/consul"
	_ "github.com/go-orb/plugins/server/http/router/chi"
)

// Injectors from wire.go:

// newComponents combines everything above and returns a slice of components.
func newComponents(serviceName types.ServiceName, serviceVersion types.ServiceVersion) ([]types.Component, error) {
	configData, err := provideConfigData(serviceName, serviceVersion)
	if err != nil {
		return nil, err
	}
	v := _wireValue
	logger, err := log.Provide(serviceName, configData, v...)
	if err != nil {
		return nil, err
	}
	v2 := _wireValue2
	registryType, err := registry.Provide(serviceName, serviceVersion, configData, logger, v2...)
	if err != nil {
		return nil, err
	}
	v3, err := provideServerOpts()
	if err != nil {
		return nil, err
	}
	serverServer, err := server.Provide(serviceName, configData, logger, registryType, v3...)
	if err != nil {
		return nil, err
	}
	v4, err := provideComponents(serviceName, serviceVersion, configData, logger, registryType, serverServer)
	if err != nil {
		return nil, err
	}
	return v4, nil
}

var (
	_wireValue  = []log.Option{}
	_wireValue2 = []registry.Option{}
)

// wire.go:

// provideConfigData reads the config from cli and returns it.
func provideConfigData(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
) (types.ConfigData, error) {
	u, err := url.Parse("cli://urfave")
	if err != nil {
		return nil, err
	}

	cfgSections := types.SplitServiceName(serviceName)

	data, err := config.Read([]*url.URL{u}, cfgSections)

	return data, err
}

// provideServerOpts provides options for the go-orb server.
func provideServerOpts() ([]server.ConfigOption, error) {

	hInstance := new(echo.Handler)
	hRegister := echo2.RegisterEchoHandler(hInstance)

	opts := []server.ConfigOption{}
	opts = append(opts, server.WithEntrypointConfig(grpc.NewConfig(grpc.WithName("grpc"), grpc.WithInsecure(), grpc.WithHandlers(hRegister))))
	opts = append(opts, server.WithEntrypointConfig(http.NewConfig(http.WithName("http"), http.WithInsecure(), http.WithHandlers(hRegister))))
	opts = append(opts, server.WithEntrypointConfig(http.NewConfig(http.WithName("https"), http.WithDisableHTTP2(), http.WithHandlers(hRegister))))
	opts = append(opts, server.WithEntrypointConfig(http.NewConfig(http.WithName("h2c"), http.WithInsecure(), http.WithAllowH2C(), http.WithHandlers(hRegister))))
	opts = append(opts, server.WithEntrypointConfig(http.NewConfig(http.WithName("http2"), http.WithHandlers(hRegister))))
	opts = append(opts, server.WithEntrypointConfig(http.NewConfig(http.WithName("http3"), http.WithHTTP3(), http.WithHandlers(hRegister))))
	opts = append(opts, server.WithEntrypointConfig(drpc.NewConfig(drpc.WithName("drpc"), drpc.WithHandlers(hRegister))))

	return opts, nil
}

// provideComponents creates a slice of components out of the arguments.
func provideComponents(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	cfgData types.ConfigData,
	logger log.Logger,
	reg registry.Type,
	srv server.Server,
) ([]types.Component, error) {
	components := []types.Component{}
	components = append(components, logger)
	components = append(components, reg)
	components = append(components, &srv)

	return components, nil
}
