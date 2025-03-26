//go:build wireinject
// +build wireinject

package service

import (
	"context"
	"fmt"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/wire"
)

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
) (ActionServer, error) {
	return func() error {
		// Orb start
		for _, c := range components.Iterate(false) {
			logger.Debug("Starting", "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))

			err := c.Start(serviceContext.Context())
			if err != nil {
				logger.Error("Failed to start", "error", err, "component", fmt.Sprintf("%s/%s", c.Type(), c.String()))
				return fmt.Errorf("failed to start component %s/%s: %w", c.Type(), c.String(), err)
			}
		}

		// Let the service work.
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

// ProvideRunner provides a runner for the service.
func ProvideRunner(
	appContext *cli.AppContext,
	appConfigData cli.AppConfigData,
	flags []*cli.Flag,
) (Runner, error) {
	panic(wire.Build(
		provideServiceContext,
		types.ProvideComponents,

		cli.ProvideServiceConfigData,

		log.ProvideWithServiceNameField,

		provideActionServer,
		provideActionHealth,
		provideRunner,
	))
}
