package main

import (
	"context"
	"os"

	"github.com/go-orb/examples/event/simple/pb/user_new"
	"github.com/go-orb/go-orb/event"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/config/source/cli/urfave"
	_ "github.com/go-orb/plugins/event/natsjs"
	_ "github.com/go-orb/plugins/log/slog"
	"github.com/google/uuid"
)

func runner(
	sn types.ServiceName,
	configs types.ConfigData,
	logger log.Logger,
	eventWire event.Handler,
	done chan os.Signal,
) error {

	userNewHandler := func(ctx context.Context, req *user_new.Request) (*user_new.Resp, error) {
		return &user_new.Resp{Name: req.GetName(), Uuid: uuid.New().String()}, nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	event.HandleRequest(ctx, eventWire, "user.new", userNewHandler)

	// Blocks until sigterm/sigkill.
	<-done

	cancel()

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
