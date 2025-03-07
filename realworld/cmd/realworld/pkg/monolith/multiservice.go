// Package monolith implements the realworld example app.
package monolith

import (
	"github.com/go-orb/go-orb/cli"
)

// Name is the name of the app.
const Name = "realworld"

// Version is the version of the app.
//
//nolint:gochecknoglobals
var Version = ""

// MainCommands returns commands specific to the service.
func MainCommands() []*cli.Command {
	return []*cli.Command{
		{
			Service: Name,
			Name:    "server",
			Usage:   "Start the servers",
		},
		{
			Service: Name,
			Name:    "health",
			Usage:   "Check the health of the services",
		},
	}
}
