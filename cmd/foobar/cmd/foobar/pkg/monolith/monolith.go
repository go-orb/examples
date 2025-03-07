package monolith

import (
	"github.com/go-orb/go-orb/cli"
)

const Name = "foobar"

var Version = "unset"

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
