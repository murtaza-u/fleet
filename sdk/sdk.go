package sdk

import (
	"context"
	"net/http"
)

// Handle connects to the gRPC Fleet server, managing the proxied
// request using the provided HTTP handler.
func Handle(h http.Handler, opts ...OptFunc) error {
	ctx := context.Background()
	return HandleWithCtx(ctx, h, opts...)
}

// HandleWithCtx is similar to Handle but supports passing a context.
func HandleWithCtx(ctx context.Context, h http.Handler, opts ...OptFunc) error {
	c, err := newClient(ctx, opts...)
	if err != nil {
		return err
	}
	return c.ListenAndServe(h)
}
