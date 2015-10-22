package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/eduardonunesp/sslb/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/eduardonunesp/sslb/cfg"
)

const APP_NAME = "sslb (SUPER SIMPLE LOAD BALANCER)"
const APP_USAGE = "sslb"
const VERSION_MAJOR = "0"
const VERSION_MINOR = "0"
const VERSION_BUILD = "3"
const CONFIG_FILENAME = "config.json.example"

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, b",
			Usage: "activate the verbose output",
		},
		cli.BoolFlag{
			Name:  "config, c",
			Usage: "create an example of config file",
		},
	}

	app.Name = APP_NAME
	app.Usage = APP_USAGE
	app.Version = VERSION_MAJOR + "." + VERSION_MINOR + "." + VERSION_BUILD

	app.Action = func(c *cli.Context) {
		if !c.Bool("verbose") {
			log.SetOutput(ioutil.Discard)
		}

		if c.Bool("config") {
			cfg.CreateConfig(CONFIG_FILENAME)
			os.Exit(0)
		}

		// The function setup do everything for configure
		// and return the server ready to run
		server := cfg.Setup()
		server.Run()

	}

	app.Run(os.Args)
}
