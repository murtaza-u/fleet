package cli

import (
	"path/filepath"

	"github.com/murtaza-u/fleet/srv"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func runFlags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewUintFlag(&cli.UintFlag{
			Name:    "rpc-port",
			Usage:   "Port for the gRPC server",
			Value:   2035,
			EnvVars: []string{"FLEET_RPC_PORT"},
		}),
		altsrc.NewUintFlag(&cli.UintFlag{
			Name:    "http-port",
			Usage:   "Port for the http web server",
			Value:   8080,
			EnvVars: []string{"FLEET_HTTP_PORT"},
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:    "reflect",
			Usage:   "Enable gRPC reflection",
			EnvVars: []string{"FLEET_REFLECT"},
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:    "tls",
			Usage:   "Enable server-side gRPC TLS",
			EnvVars: []string{"FLEET_TLS"},
		}),
		altsrc.NewPathFlag(&cli.PathFlag{
			Name:      "certificate",
			Usage:     "Path to certificate file",
			Value:     filepath.Join("certs", "srv-cert.pem"),
			TakesFile: true,
			EnvVars:   []string{"FLEET_CERTIFICATE_PATH"},
		}),
		altsrc.NewPathFlag(&cli.PathFlag{
			Name:      "private-key",
			Usage:     "Path to private key file",
			Value:     filepath.Join("certs", "srv-key.pem"),
			TakesFile: true,
			EnvVars:   []string{"FLEET_PRIVATE_KEY_FILE"},
		}),
		altsrc.NewInt64Flag(&cli.Int64Flag{
			Name:    "max-http-request-body-size",
			Usage:   "Maximum size of the HTTP request body",
			Value:   1024 * 1024,
			EnvVars: []string{"FLEET_MAX_HTTP_REQUEST_BODY_SIZE"},
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:    "max-grpc-recv-size",
			Usage:   "Maximum gRPC receive message size",
			Value:   1024 * 1024,
			EnvVars: []string{"FLEET_MAX_GRPC_RECV_SIZE"},
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:    "max-grpc-send-size",
			Usage:   "Maximum gRPC send message size.",
			Value:   1024 * 1024,
			EnvVars: []string{"FLEET_MAX_GRPC_SEND_SIZE"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "serving-url-format",
			Usage:   "Serving url format",
			Value:   "http://%s.localhost:8080",
			EnvVars: []string{"FLEET_SERVING_URL_FORMAT"},
		}),
		&cli.PathFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to config file",
			Value:   filepath.Join(".", "config.yaml"),
			EnvVars: []string{"FLEET_CONFIG_FILE"},
		},
	}
}

func runCmd() *cli.Command {
	flags := runFlags()

	return &cli.Command{
		Name:   "run",
		Usage:  "Runs the fleet server",
		Args:   false,
		Flags:  flags,
		Before: loadConfigIfExists(flags),
		Action: func(ctx *cli.Context) error {
			opts := []srv.OptFunc{
				srv.WithRPCPort(uint16(ctx.Uint("rpc-port"))),
				srv.WithHTTPPort(uint16(ctx.Uint("http-port"))),
				srv.WithPathToCertificate(ctx.Path("certificate")),
				srv.WithPathToPrivateKey(ctx.Path("private-key")),
				srv.WithMaxRequestBodySize(ctx.Int64("max-http-request-body-size")),
				srv.WithMaxMsgRecvSize(ctx.Int("max-grpc-recv-size")),
				srv.WithMaxMsgRecvSize(ctx.Int("max-grpc-send-size")),
				srv.WithServingUrlFormat(ctx.String("serving-url-format")),
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
}
