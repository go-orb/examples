// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/go-orb/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package service

import (
	"context"
	"fmt"
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
)

// Injectors from wire.go:

// ProvideRunner provides a runner for the service.
func ProvideRunner(appContext *cli.AppContext, flags []*cli.Flag) (Runner, error) {
	serviceContext, err := provideServiceContext(appContext)
	if err != nil {
		return nil, err
	}
	v, err := types.ProvideComponents()
	if err != nil {
		return nil, err
	}
	serviceName, err := cli.ProvideServiceName(serviceContext)
	if err != nil {
		return nil, err
	}
	configData, err := cli.ProvideConfigData(serviceContext, flags)
	if err != nil {
		return nil, err
	}
	v2 := _wireValue
	logger, err := log.Provide(serviceName, configData, v, v2...)
	if err != nil {
		return nil, err
	}
	actionServer, err := provideActionServer(serviceContext, v, logger)
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

var (
	_wireValue = []log.Option{}
)

// wire.go:

type Runner func(command []string) error

type ActionServer func() error

type ActionHealth func() error

func provideServiceContext(appContext *cli.AppContext) (*cli.ServiceContext, error) {
	return cli.NewServiceContext(appContext, Name, Version), nil
}

func provideActionServer(
	serviceContext *cli.ServiceContext,
	components *types.Components,
	logger log.Logger,
) (ActionServer, error) {
	return func() error {

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
