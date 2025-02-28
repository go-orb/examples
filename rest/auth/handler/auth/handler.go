// Package echo provdes a echo handler.
package echo

import (
	"context"
	"fmt"
	"strings"
	"time"

	authV1Proto "github.com/go-orb/examples/rest/auth/proto/auth_v1"
	"github.com/go-orb/go-orb/log"
	"github.com/go-orb/go-orb/util/metadata"
	"github.com/go-orb/go-orb/util/orberrors"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/protobuf/types/known/emptypb"
)

const authorizationHeaderName = "authorization"
const authorizationPrefix = "Bearer "

var _ authV1Proto.AuthHandler = (*Handler)(nil)

// Handler is a test handler.
type Handler struct {
	tokenSecret []byte
	logger      log.Logger
}

// New creates a new Handler.
func New(tokenSecret []byte, logger log.Logger) *Handler {
	return &Handler{
		tokenSecret: tokenSecret,
		logger:      logger,
	}
}

// Login implements the fake login method.
func (h *Handler) Login(_ context.Context, req *authV1Proto.LoginRequest) (*authV1Proto.LoginResponse, error) {
	if req.GetUsername() != "someUserName" || req.GetPassword() != "changeMe" {
		return nil, orberrors.ErrUnauthorized
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": req.GetUsername(),
		"nbf": time.Now().Unix(),
	})

	tokenString, err := token.SignedString(h.tokenSecret)
	if err != nil {
		h.logger.Error("while signing token", "error", err)
		return nil, orberrors.ErrInternalServerError
	}

	resp := &authV1Proto.LoginResponse{
		Token: tokenString,
	}

	return resp, nil
}

// Introspect returns data from the token.
func (h *Handler) Introspect(ctx context.Context, _ *emptypb.Empty) (*authV1Proto.IntrospectResponse, error) {
	md, ok := metadata.Incoming(ctx)
	if !ok {
		return nil, orberrors.ErrUnauthorized
	}

	if _, ok2 := md[authorizationHeaderName]; !ok2 || !strings.HasPrefix(md[authorizationHeaderName], authorizationPrefix) {
		return nil, orberrors.ErrUnauthorized
	}

	token, err := jwt.Parse(md[authorizationHeaderName][7:], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			h.logger.Error("unexpected signing method", "alg", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return h.tokenSecret, nil
	})
	if err != nil {
		h.logger.Error("while reading a token", "error", err)
		return nil, orberrors.ErrInternalServerError
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		h.logger.Error("while parsing claims")
		return nil, orberrors.ErrInternalServerError
	}

	return &authV1Proto.IntrospectResponse{Username: claims["sub"].(string)}, nil //nolint:errcheck
}
