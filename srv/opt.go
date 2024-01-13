package srv

import (
	"fmt"
	"os"
	"path/filepath"
)

type OptFunc func(*opts)

type opts struct {
	// port for the gRPC server.
	rpcPort string
	// port for the http web server.
	httpPort string
	// enable gRPC reflection
	reflect bool
	// log in JSON
	jsonLogger bool
	// enable debug logs
	debug bool
	// enable server-side TLS
	tls bool
	// path to certificate file
	cert string
	// path to private key file
	priv string
	// limit the HTTP request body size to prevent potential abuse.
	maxRequestBodySize int64
	// maximum gRPC receive message size
	maxRecvSize int
	// maximum gRPC send message size
	maxSendSize int
}

// WithRPCPort sets the port for the gRPC server. Ensure that the chosen
// port is available and not in use.
//
// Default: 2035
func WithRPCPort(port uint16) OptFunc {
	return func(o *opts) {
		o.rpcPort = fmt.Sprintf(":%d", port)
	}
}

// WithHTTPPort sets the port for the http web server. Ensure that the
// chosen port is available and not in use.
//
// Default: 8080
func WithHTTPPort(port uint16) OptFunc {
	return func(o *opts) {
		o.httpPort = fmt.Sprintf(":%d", port)
	}
}

// WithReflection enables gRPC reflection
//
// Default: disabled
func WithReflection() OptFunc {
	return func(o *opts) {
		o.reflect = true
	}
}

// WithJSONLogger configures logger to use JSON.
//
// Default: disabled
func WithJSONLogger() OptFunc {
	return func(o *opts) {
		o.jsonLogger = true
	}
}

// WithDebug enables debug logs.
//
// Default: disabled
func WithDebug() OptFunc {
	return func(o *opts) {
		o.debug = true
	}
}

// WithPathToCertificate sets the path to the TLS certficiate file.
//
// Defaullt: ./certs/srv-cert.pem
func WithPathToCertificate(path string) OptFunc {
	return func(o *opts) {
		o.cert = path
	}
}

// WithPathToPrivateKey sets the path to the TLS private key file.
//
// Defaullt: ./certs/srv-key.pem
func WithPathToPrivateKey(path string) OptFunc {
	return func(o *opts) {
		o.priv = path
	}
}

// WithTLS enables server-side TLS.
//
// Default: false
func WithTLS() OptFunc {
	return func(o *opts) {
		o.tls = true
	}
}

// WithMaxRequestBodySize sets a limit on the size of the HTTP request
// body to prevent potential abuse.
//
// Default: 1 MB (1024 * 1024)
func WithMaxRequestBodySize(size int64) OptFunc {
	return func(o *opts) {
		o.maxRequestBodySize = size
	}
}

// WithMaxMsgRecvSize sets the maximum gRPC receive message size.
//
// Default: 1 MB (1024 * 1024)
func WithMaxMsgRecvSize(size int) OptFunc {
	return func(o *opts) {
		o.maxRecvSize = size
	}
}

// WithMaxMsgSendSize sets the maximum gRPC send message size.
//
// Default: 1 MB (1024 * 1024)
func WithMaxMsgSendSize(size int) OptFunc {
	return func(o *opts) {
		o.maxSendSize = size
	}
}

func (o opts) validate() error {
	if !o.tls {
		return nil
	}
	if _, err := os.Stat(o.cert); err != nil {
		return fmt.Errorf(
			"failed to get statistics for path %q: %w", o.cert, err)
	}
	if _, err := os.Stat(o.priv); err != nil {
		return fmt.Errorf(
			"failed to get statistics for path %q: %w", o.priv, err)
	}
	return nil
}

func defaultOpts() opts {
	return opts{
		rpcPort:            ":2035",
		httpPort:           ":8080",
		reflect:            false,
		jsonLogger:         false,
		debug:              false,
		cert:               filepath.Join("certs", "srv-cert.pem"),
		priv:               filepath.Join("certs", "srv-key.pem"),
		tls:                false,
		maxRequestBodySize: 1024 * 1024,
		maxSendSize:        1024 * 1024,
		maxRecvSize:        1024 * 1024,
	}
}
