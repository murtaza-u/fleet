package srv

import (
	"fmt"
	"os"
	"path/filepath"
)

type OptFunc func(*opts) error

type opts struct {
	// port for the gRPC server.
	rpcPort string
	// port for the http web server.
	httpPort string
	// limit the request body size to prevent potential abuse.
	maxRequestBodySize int64
	// enable gRPC reflection
	reflect bool
	// log in JSON
	jsonLogger bool
	// enable debug logs
	debug bool
	// path to certificate file
	cert string
	// path to private key file
	priv string
	// enable server-side TLS
	tls bool
}

// WithRPCPort sets the port for the gRPC server. Ensure that the chosen
// port is available and not in use.
//
// Default: 2035
func WithRPCPort(port uint16) OptFunc {
	return func(o *opts) error {
		o.rpcPort = fmt.Sprintf(":%d", port)
		return nil
	}
}

// WithHTTPPort sets the port for the http web server. Ensure that the
// chosen port is available and not in use.
//
// Default: 8080
func WithHTTPPort(port uint16) OptFunc {
	return func(o *opts) error {
		o.httpPort = fmt.Sprintf(":%d", port)
		return nil
	}
}

// WithMaxRequestBodySize sets a limit on the size of the request body
// to prevent potential abuse.
//
// Default: 1 MB
func WithMaxRequestBodySize(size int64) OptFunc {
	return func(o *opts) error {
		o.maxRequestBodySize = size
		return nil
	}
}

// WithReflection enables gRPC reflection
//
// Default: disabled
func WithReflection() OptFunc {
	return func(o *opts) error {
		o.reflect = true
		return nil
	}
}

// WithJSONLogger configures logger to use JSON.
//
// Default: disabled
func WithJSONLogger() OptFunc {
	return func(o *opts) error {
		o.jsonLogger = true
		return nil
	}
}

// WithDebug enables debug logs.
//
// Default: disabled
func WithDebug() OptFunc {
	return func(o *opts) error {
		o.debug = true
		return nil
	}
}

// WithPathToCertificate sets the path to the TLS certficiate file.
//
// Defaullt: ./certs/srv-cert.pem
func WithPathToCertificate(path string) OptFunc {
	return func(o *opts) error {
		_, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to get statistics for path %q: %w",
				path, err)
		}
		o.cert = path
		return nil
	}
}

// WithPathToPrivateKey sets the path to the TLS private key file.
//
// Defaullt: ./certs/srv-key.pem
func WithPathToPrivateKey(path string) OptFunc {
	return func(o *opts) error {
		_, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to get statistics for path %q: %w",
				path, err)
		}
		o.priv = path
		return nil
	}
}

// WithTLS enables server-side TLS.
//
// Default: false
func WithTLS() OptFunc {
	return func(o *opts) error {
		o.tls = true
		return nil
	}
}

func defaultOpts() opts {
	return opts{
		rpcPort:            ":2035",
		httpPort:           ":8080",
		maxRequestBodySize: 10 ^ 6,
		reflect:            false,
		jsonLogger:         false,
		debug:              false,
		cert:               filepath.Join("certs", "srv-cert.pem"),
		priv:               filepath.Join("certs", "srv-key.pem"),
		tls:                false,
	}
}
