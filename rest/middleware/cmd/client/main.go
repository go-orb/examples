// Package main contains a go-orb client with utilizes middlewares.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/log"

	// Own imports.
	echoproto "github.com/go-orb/examples/rest/middleware/proto/echo"

	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/slog"

	_ "github.com/go-orb/plugins/registry/consul"
	_ "github.com/go-orb/plugins/registry/mdns"

	_ "github.com/go-orb/plugins/client/middleware/log"
	_ "github.com/go-orb/plugins/client/orb"
	_ "github.com/go-orb/plugins/client/orb_transport/drpc"
	_ "github.com/go-orb/plugins/client/orb_transport/grpc"
)

func runner(
	ctx context.Context,
	logger log.Logger,
	clientWire client.Type,
) error {
	// Create a request.
	req := &echoproto.Req{Payload: []byte("Hello World")}

	// Run the query.
	protoClient := echoproto.NewEchoClient(clientWire)
	resp, err := protoClient.Echo(ctx, "orb.examples.rest.middleware.server", req)

	if err != nil {
		logger.Error("while requesting", "error", err)
		return err
	}

	if !slices.Equal(resp.GetPayload(), req.GetPayload()) {
		logger.Error("while requesting", "expected", req.GetPayload(), "got", resp.GetPayload())
		return errors.New("bad response")
	}

	return nil
}

func main() {
	app := cli.App{
		Name:     "orb.examples.rest.middleware.client",
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
