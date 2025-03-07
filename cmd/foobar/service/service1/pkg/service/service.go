package service

import (
	"fmt"

	"github.com/go-orb/go-orb/cli"
)

const Name = "foobar.service1"

var Version = "unset"

// MainCommands returns the commands which get appended to the "main/monolith" App.
func MainCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "service1",
			Service:     Name,
			Category:    "service",
			Subcommands: ServiceCommands(),
			NoAction:    true,
		},
	}
}

// ServiceCommands returns commands specific to the service.
func ServiceCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "server",
			Service: Name,
			Usage:   fmt.Sprintf("Start the %s server", Name),
		},
		{
			Name:    "health",
			Service: Name,
			Usage:   fmt.Sprintf("Check the health of the %s service", Name),
		},
	}
}
