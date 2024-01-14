package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func loadConfigIfExists(flags []cli.Flag) cli.BeforeFunc {
	def := altsrc.NewMapInputSource("", make(map[interface{}]interface{}))
	return altsrc.InitInputSourceWithContext(
		flags,
		func(ctx *cli.Context) (altsrc.InputSourceContext, error) {
			path := ctx.Path("config")
			stat, err := os.Stat(path)
			if err != nil {
				return def, nil
			}
			if stat.IsDir() {
				return def, fmt.Errorf("provided config %q is a directory",
					path)
			}
			return altsrc.NewYamlSourceFromFile(path)
		},
	)
}
