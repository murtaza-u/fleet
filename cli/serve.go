package cli

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/murtaza-u/fleet/sdk"

	"github.com/lucasepe/codename"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func serveFlags() []cli.Flag {
	conf, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(conf, "fleet", "config.yaml")

	return []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "address",
			Aliases: []string{"addr"},
			Usage:   "Address(host:port) of the Fleet gRPC server",
			Value:   "localhost:2035",
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:  "tls",
			Usage: "Use TLS to connect to the Fleet gRPC server",
		}),
		altsrc.NewPathFlag(&cli.PathFlag{
			Name:      "ca-cert",
			Usage:     "Path to the CA's certificate file",
			Value:     filepath.Join("certs", "ca-cert.pem"),
			TakesFile: true,
		}),
		&cli.StringFlag{
			Name:    "propose-subdomain",
			Aliases: []string{"p"},
			Usage:   "Propose a subdomain to register yourself for the current session",
		},
		&cli.PathFlag{
			Name:      "config",
			Aliases:   []string{"c"},
			Usage:     "Path to config file",
			Value:     path,
			TakesFile: true,
		},
	}
}

func serveCmd() *cli.Command {
	flags := serveFlags()

	return &cli.Command{
		Name:        "serve",
		Usage:       "Connect to the fleet server and serve incoming requests",
		UsageText:   "[options] static|proxy",
		Subcommands: []*cli.Command{staticCmd, proxyCmd},
		Flags:       flags,
		Before:      loadConfigIfExists(flags),
	}
}

var staticCmd = &cli.Command{
	Name:      "static",
	Usage:     "Serves files from the provided directory",
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
	Usage:     "Proxies incoming requests to the provided URL",
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
