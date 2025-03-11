// Package main contains a simple client example for a event run.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-orb/examples/event/simple/pb/user_new"
	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/event"
	"github.com/go-orb/go-orb/log"
	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/event/natsjs"
	_ "github.com/go-orb/plugins/log/slog"
)

func runner(
	_ context.Context,
	logger log.Logger,
	eventHandler event.Type,
) error {
	req := &user_new.Request{Name: "René"}

	resp, err := event.Request[user_new.Resp](context.Background(), eventHandler, "user.new", req)
	if err != nil {
		logger.Error("while requesting", "err", err)
		return fmt.Errorf("while requesting: %w", err)
	}

	logger.Info("New id for user", "name", resp.GetName(), "uuid", resp.GetUuid())

	return nil
}

func main() {
	app := cli.App{
		Name:     "orb.examples.event.simple.request",
		Version:  "",
		Usage:    "A foobar example app",
		NoAction: false,
		Flags: []*cli.Flag{
			{
				Name:        "log_level",
				Default:     "INFO",
				EnvVars:     []string{"LOG_LEVEL"},
				ConfigPaths: []cli.FlagConfigPath{{Path: []string{"logger", "level"}}},
				Usage:       "Set the log level, one of TRACE, DEBUG, INFO, WARN, ERROR",
			},
		},
		Commands: []*cli.Command{},
	}

	appContext := cli.NewAppContext(&app)

	_, err := run(appContext, os.Args, runner)
	if err != nil {
		//nolint:forbidigo
		fmt.Printf("run error: %s\n", err)
		os.Exit(1)
	}
}
