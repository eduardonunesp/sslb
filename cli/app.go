package cli

import (
	"os"

	"github.com/eduardonunesp/sslb/Godeps/_workspace/src/github.com/codegangsta/cli"
)

const APP_NAME = "SSLB (github.com/eduardonunesp/sslb)"
const APP_USAGE = "sslb"
const VERSION_MAJOR = "0"
const VERSION_MINOR = "1"
const VERSION_BUILD = "0"

func CreateAPP() *cli.App {
	app := cli.NewApp()

	app.Name = APP_NAME
	app.Usage = APP_USAGE
	app.Version = VERSION_MAJOR + "." + VERSION_MINOR + "." + VERSION_BUILD

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, b",
			Usage: "activate the verbose output",
		},
		cli.StringFlag{
			Name:  "filename, f",
			Usage: "set the filename as the configuration",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "Return the internal status",
			Action:  InternalStatus,
		},
	}

	app.Action = RunServer
	app.Run(os.Args)

	return app
}
