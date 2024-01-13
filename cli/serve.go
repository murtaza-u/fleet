package cli

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/murtaza-u/fleet/sdk"

	"github.com/lucasepe/codename"
	"github.com/urfave/cli/v2"
)

var serveCmd = &cli.Command{
	Name:        "serve",
	Usage:       "connect to the fleet server and serve incoming requests",
	UsageText:   "[options] static|proxy",
	Subcommands: []*cli.Command{staticCmd, proxyCmd},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "tls",
			Usage: "Use TLS to connect to the Fleet gRPC server",
		},
		&cli.StringFlag{
			Name:    "address",
			Aliases: []string{"addr", "rpc-addr", "rpc-address"},
			Usage:   "address(host:port) of the Fleet gRPC server",
			Value:   "localhost:2035",
		},
		&cli.PathFlag{
			Name:      "ca-cert",
			Value:     filepath.Join("certs", "ca-cert.pem"),
			TakesFile: true,
		},
		&cli.StringFlag{
			Name:    "propose-subdomain",
			Aliases: []string{"subdomain", "propose", "p"},
			Usage:   "propose a subdomain to register yourself for the current session",
		},
	},
}

var staticCmd = &cli.Command{
	Name:      "static",
	Usage:     "serves files from the provided directory",
	UsageText: "[path]",
	ArgsUsage: "[path]",
	Args:      true,
	Action: func(ctx *cli.Context) error {
		path := ctx.Args().First()
		if path == "" {
			path = "."
		}
		stat, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("%q doesn't exist", path)
			}
			return err
		}
		if !stat.IsDir() {
			return fmt.Errorf("%q is not a directory", path)
		}

		opts, err := sdkOpts(ctx.Lineage()[0])
		if err != nil {
			return err
		}

		mux := http.NewServeMux()
		mux.Handle("/", http.FileServer(http.Dir(path)))

		return sdk.Handle(mux, opts...)
	},
}

var proxyCmd = &cli.Command{
	Name:      "proxy",
	Usage:     "proxies incoming requests to the provided URL",
	UsageText: "URL",
	ArgsUsage: "URL",
	Args:      true,
	Action: func(ctx *cli.Context) error {
		raw := ctx.Args().First()
		if raw == "" {
			return fmt.Errorf("missing url for proxy")
		}
		url, err := url.Parse(raw)
		if err != nil {
			return fmt.Errorf("failed to parse url %q: %w", raw, err)
		}

		opts, err := sdkOpts(ctx.Lineage()[0])
		if err != nil {
			return err
		}

		mux := http.NewServeMux()
		proxy := httputil.NewSingleHostReverseProxy(url)
		mux.Handle("/", proxy)

		return sdk.Handle(mux, opts...)
	},
}

func sdkOpts(ctx *cli.Context) ([]sdk.OptFunc, error) {
	opts := []sdk.OptFunc{
		sdk.WithRPCAddress(ctx.String("address")),
		sdk.WithPathToCACertificate(ctx.Path("ca-cert")),
	}
	if ctx.Bool("tls") {
		opts = append(opts, sdk.WithTLS())
	}
	subdomain := ctx.String("propose-subdomain")
	if subdomain == "" {
		rng, err := codename.DefaultRNG()
		if err != nil {
			return nil, fmt.Errorf("failed to generate random name: %w", err)
		}
		subdomain = codename.Generate(rng, 0)
	}
	opts = append(opts, sdk.WithPreferredSubdomain(subdomain))
	return opts, nil
}
