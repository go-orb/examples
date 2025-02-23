// Package main contains a go-orb client with utilizes middlewares.
package main

import (
	"context"
	"errors"
	"slices"

	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/types"

	// Own imports.
	echoproto "github.com/go-orb/examples/rest/middleware/proto/echo"

	_ "github.com/go-orb/plugins/codecs/json"
	_ "github.com/go-orb/plugins/codecs/proto"
	_ "github.com/go-orb/plugins/codecs/yaml"
	_ "github.com/go-orb/plugins/config/source/cli/urfave"
	_ "github.com/go-orb/plugins/config/source/file"
	_ "github.com/go-orb/plugins/log/slog"

	_ "github.com/go-orb/plugins-experimental/registry/mdns"
	_ "github.com/go-orb/plugins/registry/consul"

	_ "github.com/go-orb/plugins/client/middleware/log"
	_ "github.com/go-orb/plugins/client/orb"
	_ "github.com/go-orb/plugins/client/orb/transport/drpc"
	_ "github.com/go-orb/plugins/client/orb/transport/grpc"
)

func runner(
	logger log.Logger,
	clientWire client.Type,
) error {
	// Create a request.
	req := &echoproto.Req{Payload: []byte("Hello World")}

	// Run the query.
	protoClient := echoproto.NewEchoClient(clientWire)
	resp, err := protoClient.Echo(context.Background(), "orb.examples.rest.middleware.server", req)

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
	var (
		serviceName    = types.ServiceName("orb.examples.rest.middleware.client")
		serviceVersion = types.ServiceVersion("v0.0.1")
	)

	if _, err := run(serviceName, serviceVersion, runner); err != nil {
		log.Error("while running", "err", err)
	}
}
