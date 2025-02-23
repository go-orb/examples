// Package echo provdes a echo handler.
package echo

import (
	"context"

	"github.com/go-orb/examples/rest/middleware/proto/echo"
	"github.com/go-orb/go-orb/log"
)

var _ echo.EchoHandler = (*Handler)(nil)

// Handler is a test handler.
type Handler struct {
	logger log.Logger
}

// New creates a new Handler.
func New(logger log.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// Echo implements the echo method.
func (c *Handler) Echo(_ context.Context, req *echo.Req) (*echo.Resp, error) {
	resp := &echo.Resp{
		Payload: req.GetPayload(),
	}

	return resp, nil
}
