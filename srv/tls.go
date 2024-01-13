package srv

import (
	"crypto/tls"

	"google.golang.org/grpc/credentials"
)

func loadTLSCreds(certF, privF string) (credentials.TransportCredentials, error) {
	crt, err := tls.LoadX509KeyPair(certF, privF)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{
		Certificates: []tls.Certificate{crt},
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(cfg), nil
}
