// Package main contains a handler which acts as server for the request client.
package main

import (
	"context"
	"os"

	"github.com/go-orb/examples/benchmarks/event/pb/echo"
	"github.com/go-orb/go-orb/event"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	_ "github.com/go-orb/plugins/codecs/goccyjson"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/config/source/cli/urfave"
	_ "github.com/go-orb/plugins/event/natsjs"
	_ "github.com/go-orb/plugins/log/slog"
)

func runner(
	serviceName types.ServiceName,
	serviceVersion types.ServiceVersion,
	logger log.Logger,
	eventWire event.Handler,
	done chan os.Signal,
) error {
	echoHandler := func(_ context.Context, req *echo.Req) (*echo.Resp, error) {
		return &echo.Resp{Payload: req.GetPayload()}, nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	event.HandleRequest(ctx, eventWire, "echo.echo", echoHandler)

	logger.Info("Started", "name", serviceName, "version", serviceVersion)

	// Blocks until sigterm/sigkill.
	<-done

	logger.Info("Stopping", "name", serviceName, "version", serviceVersion)

	cancel()

	return nil
}

func main() {
	var (
		serviceName    = types.ServiceName("benchmarks.event.handler")
		serviceVersion = types.ServiceVersion("v0.0.1")
	)

	if _, err := run(serviceName, serviceVersion, runner); err != nil {
		log.Error("while running", "err", err)
	}
}
