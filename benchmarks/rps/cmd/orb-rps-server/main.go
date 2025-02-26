// Package main contains a server for running tests on.
package main

import (
	"os"

	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"

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

func runner(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	logger log.Logger,
	done chan os.Signal,
) error {
	logger.Info("Started", "name", serviceName, "version", serviceVersion)

	// Blocks until the process receives a signal.
	<-done

	logger.Info("Stopping", "name", serviceName, "version", serviceVersion)

	return nil
}

func main() {
	var (
		serviceName    = types.ServiceName("benchmarks.rps.server")
		serviceVersion = types.ServiceVersion("v0.0.1")
	)

	if _, err := run(serviceName, serviceVersion, runner); err != nil {
		log.Error("while running", "err", err)
	}
}
