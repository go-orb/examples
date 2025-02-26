// Package main contains a go-orb client which uses a fake login server.
package main

import (
	"context"
	"errors"

	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"
	"github.com/go-orb/go-orb/util/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	// Own imports.
	authproto "github.com/go-orb/examples/rest/auth/proto/auth"

	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/log/slog"

	_ "github.com/go-orb/plugins-experimental/registry/mdns"

	_ "github.com/go-orb/plugins/client/middleware/log"
	_ "github.com/go-orb/plugins/client/orb"
	_ "github.com/go-orb/plugins/client/orb/transport/drpc"
)

// provideLoggerOpts returns the logger options.
func provideLoggerOpts() ([]log.Option, error) {
	return []log.Option{log.WithLevel("TRACE")}, nil
}

//nolint:unparam
func provideClientOpts() ([]client.Option, error) {
	return []client.Option{client.WithClientMiddleware(client.MiddlewareConfig{Name: "log"})}, nil
}

func runner(
	logger log.Logger,
	clientWire client.Type,
) error {
	// Create a request.
	req := &authproto.LoginRequest{Username: "someUserName", Password: "changeMe"}

	// Run the query.
	authClient := authproto.NewAuthClient(clientWire)
	tokenResp, err := authClient.Login(context.Background(), "orb.examples.rest.auth.server", req)

	if err != nil {
		logger.Error("while requesting", "error", err)
		return err
	}

	ctx, md := metadata.WithOutgoing(context.Background())
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
	var (
		serviceName    = types.ServiceName("orb.examples.rest.auth.client")
		serviceVersion = types.ServiceVersion("v0.0.1")
	)

	if _, err := run(serviceName, serviceVersion, runner); err != nil {
		log.Error("while running", "err", err)
	}
}
