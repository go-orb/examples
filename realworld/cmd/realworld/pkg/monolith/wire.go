//go:build wireinject
// +build wireinject

package monolith

import (
	"fmt"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/wire"

	httpgatewayserviceproxy "github.com/go-orb/examples/realworld/service/httpgateway/pkg/serviceproxy"
	lobbyservice "github.com/go-orb/examples/realworld/service/lobby/pkg/service"
	httpgatewayservice "github.com/go-orb/service/httpgateway/pkg/service"
)

// Runner is wire type for the Runner function.
type Runner func() error

// ActionServer is wire type for the ActionServer function.
type ActionServer func() error

func provideServiceContext(appContext *cli.AppContext) (*cli.ServiceContext, error) {
	return cli.NewServiceContext(appContext, Name, Version), nil
}

// provideActionServer provides the command action "server".
func provideActionServer(
	serviceContext *cli.ServiceContext,
	httpGatewayRunner httpgatewayservice.Runner,
	lobbyserviceRunner lobbyservice.Runner,
) (ActionServer, error) {
	return func() error {
		serviceContext.StopWaitGroup().Add(1)

		go func() {
			if err := httpGatewayRunner([]string{"server"}); err != nil {
				serviceContext.ExitAppGracefully(1)
			}
		}()

		serviceContext.StopWaitGroup().Add(1)

		go func() {
			if err := lobbyserviceRunner([]string{"server"}); err != nil {
				serviceContext.ExitAppGracefully(1)
			}
		}()

		serviceContext.StopWaitGroup().Wait()

		return nil
	}, nil
}

// provideRunner provides the runner.
func provideRunner(appContext *cli.AppContext,
	httpGatewayRunner httpgatewayservice.Runner,
	lobbyserviceRunner lobbyservice.Runner,
	actionServer ActionServer,
) (Runner, error) {
	return func() error {
		switch appContext.SelectedService {
		case httpgatewayserviceproxy.Name:
			appContext.StopWaitGroup.Add(1)
			return httpGatewayRunner(appContext.SelectedCommand[1:])
		case lobbyservice.Name:
			appContext.StopWaitGroup.Add(1)
			return lobbyserviceRunner(appContext.SelectedCommand[1:])
		case Name:
			switch appContext.SelectedCommand[0] {
			case "server":
				return actionServer()
			default:
				return fmt.Errorf("unknown action: %s", appContext.SelectedCommand[0])
			}
		default:
			return fmt.Errorf("unknown service: %s", appContext.SelectedService)
		}
	}, nil
}

// ProvideRunner provides a runner for the service.
func ProvideRunner(
	appContext *cli.AppContext,
	flags []*cli.Flag,
) (Runner, error) {
	panic(wire.Build(
		provideServiceContext,

		httpgatewayserviceproxy.ProvideRunner,
		lobbyservice.ProvideRunner,

		provideActionServer,
		provideRunner,
	))
}
