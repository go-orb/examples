package main

import (
	"context"
	"os"

	"github.com/go-orb/examples/benchmarks/event/pb/echo"
	"github.com/go-orb/go-orb/event"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/config/source/cli/urfave"
	_ "github.com/go-orb/plugins/event/natsjs"
	_ "github.com/go-orb/plugins/log/slog"
)

func runner(
	sn types.ServiceName,
	configs types.ConfigData,
	logger log.Logger,
	eventWire event.Type,
	done chan os.Signal,
) error {

	ctx := context.Background()

	echoHandler := func(ctx context.Context, req *echo.Req) (*echo.Resp, error) {
		return &echo.Resp{Payload: req.GetPayload()}, nil
	}

	cancelFunc, err := event.HandleRequest(ctx, eventWire, "echo.echo", echoHandler)
	if err != nil {
		os.Exit(1)
	}
	defer cancelFunc()

	// Blocks until sigterm/sigkill.
	<-done

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
