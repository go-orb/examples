// Package service implements the foobar.service1 service.
package service

import (
	"fmt"

	"github.com/go-orb/go-orb/cli"
)

// Name is the service name.
const Name = "realworld.service.lobby"

// Version is the service version.
//
//nolint:gochecknoglobals
var Version = ""

// MainCommands returns the commands which get appended to the "main/monolith" App.
func MainCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "lobby",
			Service:     Name,
			Category:    "service",
			Subcommands: Commands(),
			NoAction:    true,
		},
	}
}

// Commands returns commands specific to the service.
func Commands() []*cli.Command {
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
