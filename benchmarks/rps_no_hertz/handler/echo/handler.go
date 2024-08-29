// Package echo provdes a echo handler.
package echo

import (
	"context"

	"github.com/go-orb/examples/benchmarks/rps_no_hertz/proto/echo"
)

var _ echo.EchoHandler = (*Handler)(nil)

// Handler is a test handler.
type Handler struct{}

// Echo implements the echo method.
func (c *Handler) Echo(_ context.Context, req *echo.Req) (*echo.Resp, error) {
	resp := &echo.Resp{
		Payload: req.GetPayload(),
	}

	return resp, nil
}
