// Package main contains a simple client example for a event run.
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
)

func runner(
	logger log.Logger,
	eventHandler event.Handler,
) error {
	req := &user_new.Request{Name: "René"}

	resp, err := event.Request[user_new.Resp](context.Background(), eventHandler, "user.new", req)
	if err != nil {
		logger.Error("while requesting", "err", err)
		os.Exit(1)
	}

	logger.Info("New id for user", "name", resp.GetName(), "uuid", resp.GetUuid())

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
