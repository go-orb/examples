// Package handler provides a lobby_v1.LobbyService handler.
package handler

import (
	"context"
	"net/http"

	lobby_v1 "github.com/go-orb/examples/realworld/service/lobby/proto/lobby_v1"
	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/server"
	"github.com/go-orb/go-orb/types"
	httpgateway "github.com/go-orb/httpgateway"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-orb/httpgateway/proto/httpgateway_v1"
)

var _ types.Component = (*Handler)(nil)
var _ lobby_v1.LobbyServiceHandler = (*Handler)(nil)

// Handler is the handler for the lobby_v1.LobbyService service.
type Handler struct {
	serviceName string
	logger      log.Logger
	client      client.Type
	hgClient    *httpgateway.Client
	server      server.Server
}

// Start assigns the handler to the server.
func (h *Handler) Start(ctx context.Context) error {
	hRegister := lobby_v1.RegisterLobbyServiceHandler(h)

	// Add our server handler to all entrypoints.
	h.server.GetEntrypoints().Range(func(_ string, entrypoint server.Entrypoint) bool {
		entrypoint.AddHandler(hRegister)

		return true
	})

	// Register the ListGames method with the httpgateway service.
	_, err := h.hgClient.AddRoutes(ctx, &httpgateway_v1.Routes{Routes: []*httpgateway_v1.Route{
		{
			HttpMethod: http.MethodGet,
			Path:       "/lobby/",
			Service:    h.serviceName,
			Method:     lobby_v1.EndpointLobbyServiceListGames,
		},
	}})
	if err != nil {
		h.logger.Warn("Failed to register with httpgateway", "error", err)
	}

	return nil
}

// Stop does nothing.
func (h *Handler) Stop(_ context.Context) error {
	return nil
}

// Type returns the type of the component.
func (h *Handler) Type() string {
	return "handler"
}

// String returns the name of the component.
func (h *Handler) String() string {
	return lobby_v1.HandlerLobbyService
}

// ListGames is the actual implementation of the ListGames method.
func (h *Handler) ListGames(_ context.Context, _ *emptypb.Empty) (*lobby_v1.ListGamesResponse, error) {
	return &lobby_v1.ListGamesResponse{Games: []*lobby_v1.Game{{Id: "1", Name: "Game 1"}, {Id: "2", Name: "Game 2"}}}, nil
}

// New creates a new Handler.
func New(serviceName string, logger log.Logger, client client.Type, hgClient *httpgateway.Client, server server.Server) *Handler {
	return &Handler{
		serviceName: serviceName,
		logger:      logger,
		client:      client,
		hgClient:    hgClient,
		server:      server,
	}
}

// Provide provides a handler to wire.
func Provide(
	sn types.ServiceName,
	logger log.Logger,
	client client.Type,
	hgClient *httpgateway.Client,
	server server.Server,
) (*Handler, error) {
	handler := New(string(sn), logger, client, hgClient, server)

	return handler, nil
}
