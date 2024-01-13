package sdk

import (
	"bytes"
	"net/http"

	"github.com/murtaza-u/fleet/internal/pb"
)

// rw satisfies the http.ResponseWriter interface.
type rw struct {
	buffer     bytes.Buffer
	statusCode int
	header     http.Header
}

func newResponseWriter() *rw {
	return &rw{
		header: make(http.Header),
	}
}

func (w *rw) Header() http.Header {
	return w.header
}

func (w *rw) Write(b []byte) (n int, err error) {
	return w.buffer.Write(b)
}

func (w *rw) WriteHeader(status int) {
	w.statusCode = status
}

func (w rw) httpToPbHeader() map[string]*pb.HeaderFields {
	m := make(map[string]*pb.HeaderFields, len(w.header))
	for k, v := range w.header {
		m[k] = &pb.HeaderFields{
			Fields: v,
		}
	}
	return m
}
