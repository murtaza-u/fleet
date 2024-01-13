package cli

import (
	"os"

	"github.com/urfave/cli/v2"
)

const version = "24.01.13"

// Run runs the cli application.
func Run() error {
	app := cli.NewApp()
	app.Usage = "Ephemeral HTTP API tunnel"
	app.UsageText = "run|serve [options] [arguments...]"
	app.Version = version
	app.EnableBashCompletion = true
	app.Copyright = "Apache-2.0"
	app.Authors = []*cli.Author{
		{Name: "Murtaza Udaipurwala", Email: "murtaza@murtazau.xyz"},
	}
	app.Commands = []*cli.Command{runCmd, serveCmd}
	return app.Run(os.Args)
}
