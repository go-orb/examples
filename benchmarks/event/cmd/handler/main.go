package main

import (
	"context"
	"os"

	"github.com/go-orb/examples/benchmarks/event/pb/echo"
	"github.com/go-orb/go-orb/event"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	_ "github.com/go-orb/plugins/codecs/jsonpb"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/cli/urfave"
	_ "github.com/go-orb/plugins/config/source/file"
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

	inChan, cancelFunc, err := event.HandleRequest[echo.Req, echo.Resp](ctx, eventWire, "echo.echo")
	if err != nil {
		os.Exit(1)
	}
	defer cancelFunc()

LOOP:
	for {
		select {
		case <-done:
			break LOOP
		case <-ctx.Done():
			break LOOP
		case req := <-inChan:
			if req.Err != nil {
				logger.Error("while handling a request", "err", req.Err)
				continue
			}

			reply := &echo.Resp{Payload: req.Data.GetPayload()}
			if err := req.Reply(reply, nil); err != nil {
				logger.Error("while sending a reply", "err", err)
			}
		}
	}

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
