package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eduardonunesp/sslb/configs"
	"github.com/eduardonunesp/sslb/impl"
	"github.com/eduardonunesp/sslb/internal"
	"github.com/eduardonunesp/sslb/types"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "sslb",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

var verboseFlag bool
var configPathFlag string

func runServer() {
	if !verboseFlag {
		log.SetOutput(ioutil.Discard)
	}

	var server types.ServerRunner = impl.NewServer()
	server.SetConfig(configs.Parse(configPathFlag))

	var frontendManagerFactory types.FrontendManagerFactoryCreator = impl.NewFrontendManagerFactory()
	var backendManagerFactory types.BackendManagerFactoryCreator = impl.NewBackendManagerFactory()
	var frontendFactory types.FrontendFactoryCreator = impl.NewFrontendFactory()
	var backendFactory types.BackendFactoryCreator = impl.NewBackendFactory()

	internal.RunServer(
		server,
		frontendManagerFactory,
		backendManagerFactory,
		frontendFactory,
		backendFactory,
	)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	server.Stop()
}

// Execute the command
func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&configPathFlag, "configuration", "f", "config.json", "Set the filename as the configuration")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
