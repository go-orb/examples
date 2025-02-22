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
	eventWire event.Handler,
	done chan os.Signal,
) error {
	echoHandler := func(_ context.Context, req *echo.Req) (*echo.Resp, error) {
		return &echo.Resp{Payload: req.GetPayload()}, nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	event.HandleRequest(ctx, eventWire, "echo.echo", echoHandler)

	// Blocks until sigterm/sigkill.
	<-done

	cancel()

	return nil
}

func main() {
	var (
		serviceName    = types.ServiceName("orb.examples.event.bench_handler")
		serviceVersion = types.ServiceVersion("v0.0.1")
	)

	if _, err := run(serviceName, serviceVersion, runner); err != nil {
		log.Error("while running", "err", err)
	}
}
