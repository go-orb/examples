// Code generated by protoc-gen-go-orb. DO NOT EDIT.
//
// version:
// - protoc-gen-go-orb        v0.0.1
// - protoc                   v5.29.2
//
// Proto source: auth_v1/auth.proto

package auth_v1

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

	mdrpc "github.com/go-orb/plugins/server/drpc"
	memory "github.com/go-orb/plugins/server/memory"
)

// HandlerAuth is the name of a service, it's here to static type/reference.
const HandlerAuth = "auth.v1.Auth"
const EndpointAuthLogin = "/auth.v1.Auth/Login"
const EndpointAuthIntrospect = "/auth.v1.Auth/Introspect"

// orbEncoding_Auth_proto is a protobuf encoder for the auth.v1.Auth service.
type orbEncoding_Auth_proto struct{}

// Marshal implements the drpc.Encoding interface.
func (orbEncoding_Auth_proto) Marshal(msg drpc.Message) ([]byte, error) {
	m, ok := msg.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("message is not a proto.Message: %T", msg)
	}
	return proto.Marshal(m)
}

// Unmarshal implements the drpc.Encoding interface.
func (orbEncoding_Auth_proto) Unmarshal(data []byte, msg drpc.Message) error {
	m, ok := msg.(proto.Message)
	if !ok {
		return fmt.Errorf("message is not a proto.Message: %T", msg)
	}
	return proto.Unmarshal(data, m)
}

// Name implements the drpc.Encoding interface.
func (orbEncoding_Auth_proto) Name() string {
	return "proto"
}

// AuthClient is the client for auth.v1.Auth
type AuthClient struct {
	client client.Client
}

// NewAuthClient creates a new client for auth.v1.Auth
func NewAuthClient(client client.Client) *AuthClient {
	return &AuthClient{client: client}
}

// Login requests Login.
func (c *AuthClient) Login(ctx context.Context, service string, req *LoginRequest, opts ...client.CallOption) (*LoginResponse, error) {
	return client.Request[LoginResponse](ctx, c.client, service, EndpointAuthLogin, req, opts...)
}

// Introspect requests Introspect.
func (c *AuthClient) Introspect(ctx context.Context, service string, req *emptypb.Empty, opts ...client.CallOption) (*IntrospectResponse, error) {
	return client.Request[IntrospectResponse](ctx, c.client, service, EndpointAuthIntrospect, req, opts...)
}

// AuthHandler is the Handler for auth.v1.Auth
type AuthHandler interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)

	Introspect(ctx context.Context, req *emptypb.Empty) (*IntrospectResponse, error)
}

// orbDRPCAuthHandler wraps a AuthHandler to implement DRPCAuthServer.
type orbDRPCAuthHandler struct {
	handler AuthHandler
}

// Login implements the DRPCAuthServer interface by adapting to the AuthHandler.
func (w *orbDRPCAuthHandler) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	return w.handler.Login(ctx, req)
}

// Introspect implements the DRPCAuthServer interface by adapting to the AuthHandler.
func (w *orbDRPCAuthHandler) Introspect(ctx context.Context, req *emptypb.Empty) (*IntrospectResponse, error) {
	return w.handler.Introspect(ctx, req)
}

// Stream adapters to convert DRPC streams to ORB streams.

// Verification that our adapters implement the required interfaces.
var _ DRPCAuthServer = (*orbDRPCAuthHandler)(nil)

// registerAuthDRPCHandler registers the service to an dRPC server.
func registerAuthDRPCHandler(srv *mdrpc.Server, handler AuthHandler) error {
	desc := DRPCAuthDescription{}

	// Wrap the ORB handler with our adapter to make it compatible with DRPC.
	drpcHandler := &orbDRPCAuthHandler{handler: handler}

	// Register with the server/drpc(.Mux).
	err := srv.Router().Register(drpcHandler, desc)
	if err != nil {
		return err
	}

	// Add each endpoint name of this handler to the orb drpc server.
	srv.AddEndpoint("/auth.v1.Auth/Login")
	srv.AddEndpoint("/auth.v1.Auth/Introspect")

	return nil
}

// registerAuthMemoryHandler registers the service to a memory server.
func registerAuthMemoryHandler(srv *memory.Server, handler AuthHandler) error {
	desc := DRPCAuthDescription{}

	// Wrap the ORB handler with our adapter to make it compatible with DRPC.
	drpcHandler := &orbDRPCAuthHandler{handler: handler}

	// Register with the server/drpc(.Mux).
	err := srv.Router().Register(drpcHandler, desc)
	if err != nil {
		return err
	}

	// Add each endpoint name of this handler to the orb drpc server.
	srv.AddEndpoint("/auth.v1.Auth/Login")
	srv.AddEndpoint("/auth.v1.Auth/Introspect")

	return nil
}

// RegisterAuthHandler will return a registration function that can be
// provided to entrypoints as a handler registration.
func RegisterAuthHandler(handler any) server.RegistrationFunc {
	return func(s any) {
		switch srv := s.(type) {

		case *mdrpc.Server:
			registerAuthDRPCHandler(srv, handler.(AuthHandler))
		case *memory.Server:
			registerAuthMemoryHandler(srv, handler.(AuthHandler))
		default:
			log.Warn("No provider for this server found", "proto", "auth_v1/auth.proto", "handler", "Auth", "server", s)
		}
	}
}
