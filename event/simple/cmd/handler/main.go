package main

import (
	"context"
	"os"

	"github.com/go-orb/examples/event/simple/pb/user_new"
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
	"github.com/google/uuid"
)

func runner(
	sn types.ServiceName,
	configs types.ConfigData,
	logger log.Logger,
	eventWire event.Type,
	done chan os.Signal,
) error {

	ctx := context.Background()

	inChan, cancelFunc, err := event.HandleRequest[user_new.Request, user_new.Resp](ctx, eventWire, "user.new")
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

			logger.Info("Got", "name", req.Data.GetName())
			reply := &user_new.Resp{Name: req.Data.GetName(), Uuid: uuid.New().String()}
			if err := req.Reply(reply, nil); err != nil {
				logger.Error("while sending a reply", "err", err)
			}
		}
	}

	return nil
}

func main() {
	var (
		serviceName    = types.ServiceName("orb.examples.event.simple")
		serviceVersion = types.ServiceVersion("v0.0.1")
	)

	if _, err := run(serviceName, serviceVersion, runner); err != nil {
		log.Error("while running", "err", err)
	}
}
