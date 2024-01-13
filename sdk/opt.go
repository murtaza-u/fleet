package sdk

import (
	"fmt"
	"os"
	"path/filepath"
)

type OptFunc func(*opts)

type opts struct {
	// address(host:port) of the gRPC server.
	rpcAddr string
	// preferred subdomain.
	subdomain string
	// path to CA certificate file
	caCert string
	// connect to gRPC server via TLS
	tls bool
}

// WithRPCAddress sets the address of the gRPC server.
//
// Default: 127.0.0.1:2305
func WithRPCAddress(addr string) OptFunc {
	return func(o *opts) {
		o.rpcAddr = addr
	}
}

// WithPreferredSubdomain requests the server to register the client on
// the given subdomain.
func WithPreferredSubdomain(name string) OptFunc {
	return func(o *opts) {
		o.subdomain = name
	}
}

// WithPathToCACertificate sets the path to the CA's TLS certficiate file.
//
// Defaullt: ./certs/ca-cert.pem
func WithPathToCACertificate(path string) OptFunc {
	return func(o *opts) {
		o.caCert = path
	}
}

func (o opts) validate() error {
	if o.subdomain == "" {
		return fmt.Errorf("missing subdomain")
	}
	_, err := os.Stat(o.caCert)
	if err != nil {
		return fmt.Errorf("failed to get statistics for path %q: %w",
			o.caCert, err)
	}
	return nil
}

// WithTLS sets the client to use TLS for connecting to the gRPC server.
//
// Default: false
func WithTLS() OptFunc {
	return func(o *opts) {
		o.tls = true
	}
}

func defaultOpts() opts {
	return opts{
		rpcAddr:   "127.0.0.1:2035",
		subdomain: "",
		caCert:    filepath.Join("certs", "ca-cert.pem"),
		tls:       false,
	}
}
