package cli

import (
	"path/filepath"

	"github.com/murtaza-u/fleet/srv"

	"github.com/urfave/cli/v2"
)

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "run the fleet server",
	Args:  false,
	Flags: []cli.Flag{
		&cli.UintFlag{
			Name:    "rpc-port",
			Usage:   "Port for the gRPC server",
			Value:   2035,
			EnvVars: []string{"FLEET_RPC_PORT"},
		},
		&cli.UintFlag{
			Name:    "http-port",
			Usage:   "Port for the http web server",
			Value:   8080,
			EnvVars: []string{"FLEET_HTTP_PORT"},
		},
		&cli.BoolFlag{
			Name:    "reflect",
			Usage:   "Enable gRPC reflection",
			EnvVars: []string{"FLEET_REFLECT"},
		},
		&cli.BoolFlag{
			Name:    "tls",
			Usage:   "Enable server-side gRPC TLS",
			EnvVars: []string{"FLEET_TLS"},
		},
		&cli.PathFlag{
			Name:      "certificate",
			Usage:     "path to certificate file",
			Value:     filepath.Join("certs", "srv-cert.pem"),
			TakesFile: true,
			EnvVars:   []string{"FLEET_CERTIFICATE_PATH"},
		},
		&cli.PathFlag{
			Name:      "private-key",
			Usage:     "path to private key file",
			Value:     filepath.Join("certs", "srv-key.pem"),
			TakesFile: true,
			EnvVars:   []string{"FLEET_PRIVATE_KEY_FILE"},
		},
		&cli.Int64Flag{
			Name:    "max-http-request-body-size",
			Usage:   "Maximum size of the HTTP request body",
			Value:   1024 * 1024,
			EnvVars: []string{"FLEET_MAX_HTTP_REQUEST_BODY_SIZE"},
		},
		&cli.IntFlag{
			Name:    "max-grpc-recv-size",
			Usage:   "Maximum gRPC receive message size",
			Value:   1024 * 1024,
			EnvVars: []string{"FLEET_MAX_GRPC_RECV_SIZE"},
		},
		&cli.IntFlag{
			Name:    "max-grpc-send-size",
			Usage:   "Maximum gRPC send message size.",
			Value:   1024 * 1024,
			EnvVars: []string{"FLEET_MAX_GRPC_SEND_SIZE"},
		},
	},
	Action: func(ctx *cli.Context) error {
		opts := []srv.OptFunc{
			srv.WithRPCPort(uint16(ctx.Uint("rpc-port"))),
			srv.WithHTTPPort(uint16(ctx.Uint("http-port"))),
			srv.WithPathToCertificate(ctx.Path("certificate")),
			srv.WithPathToPrivateKey(ctx.Path("private-key")),
			srv.WithMaxRequestBodySize(ctx.Int64("max-http-request-body-size")),
			srv.WithMaxMsgRecvSize(ctx.Int("max-grpc-recv-size")),
			srv.WithMaxMsgRecvSize(ctx.Int("max-grpc-send-size")),
		}
		if ctx.Bool("tls") {
			opts = append(opts, srv.WithTLS())
		}
		if ctx.Bool("reflect") {
			opts = append(opts, srv.WithReflection())
		}

		srv, err := srv.New(opts...)
		if err != nil {
			return err
		}
		return srv.Run()
	},
}
