// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/go-orb/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package service

import (
	"context"
	"fmt"
	"github.com/go-orb/examples/realworld/service/httpgateway/pkg/serviceproxy"
	"github.com/go-orb/examples/realworld/service/lobby/pkg/handler"
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/registry"
	"github.com/go-orb/go-orb/server"
	"github.com/go-orb/go-orb/types"
)

// Injectors from wire.go:

// ProvideRunner provides a runner for the service.
func ProvideRunner(appContext *cli.AppContext, appConfigData cli.AppConfigData, flags []*cli.Flag) (Runner, error) {
	serviceContext, err := provideServiceContext(appContext)
	if err != nil {
		return nil, err
	}
	serviceContextWithConfig, err := cli.ProvideServiceConfigData(serviceContext, appConfigData, flags)
	if err != nil {
		return nil, err
	}
	v, err := types.ProvideComponents()
	if err != nil {
		return nil, err
	}
	logger, err := log.ProvideWithServiceNameField(serviceContextWithConfig, v)
	if err != nil {
		return nil, err
	}
	registryType, err := registry.ProvideNoOpts(serviceContextWithConfig, v, logger)
	if err != nil {
		return nil, err
	}
	serverServer, err := server.ProvideNoOpts(serviceContextWithConfig, v, logger, registryType)
	if err != nil {
		return nil, err
	}
	clientType, err := client.ProvideNoOpts(serviceContextWithConfig, v, logger, registryType)
	if err != nil {
		return nil, err
	}
	serviceClient, err := serviceproxy.ProvideClient(clientType)
	if err != nil {
		return nil, err
	}
	handlerHandler, err := handler.Provide(serviceContext, logger, clientType, serviceClient, serverServer)
	if err != nil {
		return nil, err
	}
	actionServer, err := provideActionServer(serviceContextWithConfig, v, logger, serverServer, handlerHandler)
	if err != nil {
		return nil, err
	}
	actionHealth, err := provideActionHealth()
	if err != nil {
		return nil, err
	}
	runner, err := provideRunner(actionServer, actionHealth)
	if err != nil {
		return nil, err
	}
	return runner, nil
}

// wire.go:

type Runner func(command []string) error

type ActionServer func() error

type ActionHealth func() error

func provideServiceContext(appContext *cli.AppContext) (*cli.ServiceContext, error) {
	return cli.NewServiceContext(appContext, Name, Version), nil
}

func provideActionServer(
	serviceContext *cli.ServiceContextWithConfig,
	components *types.Components,
	logger log.Logger,
	_ server.Server, handler2 *handler.Handler,
) (ActionServer, error) {
	return func() error {
		if err := components.Add(handler2, types.PriorityHandler); err != nil {
			return err
		}

		for _, c := range components.Iterate(false) {
			logger.Debug("Starting", "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))

			err := c.Start(serviceContext.Context())
			if err != nil {
				logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
				return fmt.Errorf("failed to start component %s/%s: %w", c.Type(), c.String(), err)
			}
		}

		<-serviceContext.Context().Done()

		ctx := context.Background()

		for _, c := range components.Iterate(true) {
			logger.Debug("Stopping", "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))

			err := c.Stop(ctx)
			if err != nil {
				logger.Error("Failed to stop", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
			}
		}

		serviceContext.StopWaitGroup().Done()

		return nil
	}, nil
}

func provideActionHealth() (ActionHealth, error) {
	return func() error {
		return nil
	}, nil
}

func provideRunner(actionServer ActionServer, actionHealth ActionHealth) (Runner, error) {
	return func(command []string) error {
		switch command[0] {
		case "server":
			return actionServer()
		case "health":
			return actionHealth()
		default:
			return fmt.Errorf("unknown action: %s", command[0])
		}
	}, nil
}
