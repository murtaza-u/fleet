package sdk

import (
	"bytes"
	"net/http"
)

// rw satisfies the http.ResponseWriter interface.
type rw struct {
	buffer     bytes.Buffer
	statusCode int
}

func (rw) Header() http.Header {
	return http.Header{}
}

func (w *rw) Write(b []byte) (n int, err error) {
	return w.buffer.Write(b)
}

func (w *rw) WriteHeader(status int) {
	w.statusCode = status
}
