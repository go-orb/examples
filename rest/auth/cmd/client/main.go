// Package main contains a go-orb client which uses a fake login server.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-orb/go-orb/cli"
	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/util/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	// Own imports.
	authv1proto "github.com/go-orb/examples/rest/auth/proto/auth_v1"

	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/slog"

	_ "github.com/go-orb/plugins-experimental/registry/mdns"
	_ "github.com/go-orb/plugins/registry/consul"

	_ "github.com/go-orb/plugins/client/middleware/log"
	_ "github.com/go-orb/plugins/client/orb"
	_ "github.com/go-orb/plugins/client/orb_transport/drpc"
)

func runner(
	svcCtx *cli.ServiceContext,
	logger log.Logger,
	clientWire client.Type,
) error {
	// Create a request.
	req := &authv1proto.LoginRequest{Username: "someUserName", Password: "changeMe"}

	// Run the query.
	authClient := authv1proto.NewAuthClient(clientWire)
	tokenResp, err := authClient.Login(svcCtx.Context(), "orb.examples.rest.auth.server", req)

	if err != nil {
		logger.Error("while requesting", "error", err)
		return err
	}

	ctx, md := metadata.WithOutgoing(svcCtx.Context())
	md["authorization"] = "Bearer " + tokenResp.GetToken()

	introspectResponse, err := authClient.Introspect(ctx, "orb.examples.rest.auth.server", &emptypb.Empty{})
	if err != nil {
		logger.Error("while requesting", "error", err)
		return err
	}

	if introspectResponse.GetUsername() != req.GetUsername() {
		logger.Error("while requesting", "expected", req.GetUsername(), "got", introspectResponse.GetUsername())
		return errors.New("bad response")
	}

	logger.Info("all good")

	return nil
}

func main() {
	app := cli.App{
		Name:     "orb.examples.rest.auth.client",
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
