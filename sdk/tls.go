package sdk

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"

	"google.golang.org/grpc/credentials"
)

func loadTLSCreds(path string) (credentials.TransportCredentials, error) {
	ca, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(ca) {
		return nil, errors.New("failed to add CA's certificate")
	}
	cfg := &tls.Config{
		RootCAs: pool,
	}
	return credentials.NewTLS(cfg), nil
}
