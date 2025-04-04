// Code generated by protoc-gen-go-orb. DO NOT EDIT.
//
// version:
// - protoc-gen-go-orb        v0.0.1
// - protoc                   v5.29.2
//
// Proto source: lobby_v1/lobby_v1.proto

package lobby_v1

import (
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

import (
	"context"
	"fmt"

	"github.com/go-orb/go-orb/client"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/server"

	"google.golang.org/protobuf/proto"
	"storj.io/drpc"

	grpc "google.golang.org/grpc"

	mdrpc "github.com/go-orb/plugins/server/drpc"
	memory "github.com/go-orb/plugins/server/memory"
)

// HandlerLobbyService is the name of a service, it's here to static type/reference.
const HandlerLobbyService = "lobby.v1.LobbyService"
const EndpointLobbyServiceListGames = "/lobby.v1.LobbyService/ListGames"

// orbEncoding_LobbyService_proto is a protobuf encoder for the lobby.v1.LobbyService service.
type orbEncoding_LobbyService_proto struct{}

// Marshal implements the drpc.Encoding interface.
func (orbEncoding_LobbyService_proto) Marshal(msg drpc.Message) ([]byte, error) {
	m, ok := msg.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("message is not a proto.Message: %T", msg)
	}
	return proto.Marshal(m)
}

// Unmarshal implements the drpc.Encoding interface.
func (orbEncoding_LobbyService_proto) Unmarshal(data []byte, msg drpc.Message) error {
	m, ok := msg.(proto.Message)
	if !ok {
		return fmt.Errorf("message is not a proto.Message: %T", msg)
	}
	return proto.Unmarshal(data, m)
}

// Name implements the drpc.Encoding interface.
func (orbEncoding_LobbyService_proto) Name() string {
	return "proto"
}

// LobbyServiceClient is the client for lobby.v1.LobbyService
type LobbyServiceClient struct {
	client client.Client
}

// NewLobbyServiceClient creates a new client for lobby.v1.LobbyService
func NewLobbyServiceClient(client client.Client) *LobbyServiceClient {
	return &LobbyServiceClient{client: client}
}

// ListGames requests ListGames.
func (c *LobbyServiceClient) ListGames(ctx context.Context, service string, req *emptypb.Empty, opts ...client.CallOption) (*ListGamesResponse, error) {
	return client.Request[ListGamesResponse](ctx, c.client, service, EndpointLobbyServiceListGames, req, opts...)
}

// LobbyServiceHandler is the Handler for lobby.v1.LobbyService
type LobbyServiceHandler interface {
	ListGames(ctx context.Context, req *emptypb.Empty) (*ListGamesResponse, error)
}

// orbGRPCLobbyService provides the adapter to convert a LobbyServiceHandler to a gRPC LobbyServiceServer.
type orbGRPCLobbyService struct {
	handler LobbyServiceHandler
}

// ListGames implements the LobbyServiceServer interface by adapting to the LobbyServiceHandler.
func (s *orbGRPCLobbyService) ListGames(ctx context.Context, req *emptypb.Empty) (*ListGamesResponse, error) {
	return s.handler.ListGames(ctx, req)
}

// Stream adapters to convert gRPC streams to ORB streams.

// Verification that our adapters implement the required interfaces.
var _ LobbyServiceServer = (*orbGRPCLobbyService)(nil)

// registerLobbyServiceGRPCServerHandler registers the service to a gRPC server.
func registerLobbyServiceGRPCServerHandler(srv grpc.ServiceRegistrar, handler LobbyServiceHandler) {
	// Create the adapter to convert from LobbyServiceHandler to LobbyServiceServer
	grpcHandler := &orbGRPCLobbyService{handler: handler}

	srv.RegisterService(&LobbyService_ServiceDesc, grpcHandler)
}

// orbDRPCLobbyServiceHandler wraps a LobbyServiceHandler to implement DRPCLobbyServiceServer.
type orbDRPCLobbyServiceHandler struct {
	handler LobbyServiceHandler
}

// ListGames implements the DRPCLobbyServiceServer interface by adapting to the LobbyServiceHandler.
func (w *orbDRPCLobbyServiceHandler) ListGames(ctx context.Context, req *emptypb.Empty) (*ListGamesResponse, error) {
	return w.handler.ListGames(ctx, req)
}

// Stream adapters to convert DRPC streams to ORB streams.

// Verification that our adapters implement the required interfaces.
var _ DRPCLobbyServiceServer = (*orbDRPCLobbyServiceHandler)(nil)

// registerLobbyServiceDRPCHandler registers the service to an dRPC server.
func registerLobbyServiceDRPCHandler(srv *mdrpc.Server, handler LobbyServiceHandler) error {
	desc := DRPCLobbyServiceDescription{}

	// Wrap the ORB handler with our adapter to make it compatible with DRPC.
	drpcHandler := &orbDRPCLobbyServiceHandler{handler: handler}

	// Register with the server/drpc(.Mux).
	err := srv.Router().Register(drpcHandler, desc)
	if err != nil {
		return err
	}

	// Add each endpoint name of this handler to the orb drpc server.
	srv.AddEndpoint("/lobby.v1.LobbyService/ListGames")

	return nil
}

// registerLobbyServiceMemoryHandler registers the service to a memory server.
func registerLobbyServiceMemoryHandler(srv *memory.Server, handler LobbyServiceHandler) error {
	desc := DRPCLobbyServiceDescription{}

	// Wrap the ORB handler with our adapter to make it compatible with DRPC.
	drpcHandler := &orbDRPCLobbyServiceHandler{handler: handler}

	// Register with the server/drpc(.Mux).
	err := srv.Router().Register(drpcHandler, desc)
	if err != nil {
		return err
	}

	// Add each endpoint name of this handler to the orb drpc server.
	srv.AddEndpoint("/lobby.v1.LobbyService/ListGames")

	return nil
}

// RegisterLobbyServiceHandler will return a registration function that can be
// provided to entrypoints as a handler registration.
func RegisterLobbyServiceHandler(handler any) server.RegistrationFunc {
	return func(s any) {
		switch srv := s.(type) {

		case grpc.ServiceRegistrar:
			registerLobbyServiceGRPCServerHandler(srv, handler.(LobbyServiceHandler))
		case *mdrpc.Server:
			registerLobbyServiceDRPCHandler(srv, handler.(LobbyServiceHandler))
		case *memory.Server:
			registerLobbyServiceMemoryHandler(srv, handler.(LobbyServiceHandler))
		default:
			log.Warn("No provider for this server found", "proto", "lobby_v1/lobby_v1.proto", "handler", "LobbyService", "server", s)
		}
	}
}
