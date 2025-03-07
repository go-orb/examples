//go:build wireinject
// +build wireinject

package monolith

import (
	"fmt"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/wire"

	service1 "github.com/go-orb/examples/cmd/foobar/service/service1/pkg/service"
)

// Runner is wire type for the Runner function.
type Runner func(command []string) error

// ActionServer is wire type for the ActionServer function.
type ActionServer func() error

func provideServiceContext(appContext *cli.AppContext) (*cli.ServiceContext, error) {
	return cli.NewServiceContext(appContext, Name, Version), nil
}

// provideActionServer provides the command action "server".
func provideActionServer(serviceContext *cli.ServiceContext, service1Runner service1.Runner) (ActionServer, error) {
	return func() error {
		serviceContext.StopWaitGroup().Add(1)

		go func() {
			if err := service1Runner([]string{"server"}); err != nil {
				serviceContext.StopWaitGroup().Done()
				serviceContext.ExitAppGracefully(1)
			}
		}()

		serviceContext.StopWaitGroup().Wait()

		return nil
	}, nil
}

// provideRunner provides the runner.
func provideRunner(serviceContext *cli.ServiceContext, actionServer ActionServer) (Runner, error) {
	return Runner(func(command []string) error {
		switch command[0] {
		case "server":
			return actionServer()
		default:
			return fmt.Errorf("unknown action: %s", command[0])
		}
	}), nil
}

// ProvideRunner provides a runner for the service.
func ProvideRunner(
	appContext *cli.AppContext,
	flags []*cli.Flag,
) (Runner, error) {
	panic(wire.Build(
		provideServiceContext,

		service1.ProvideRunner,

		provideActionServer,
		provideRunner,
	))
}
