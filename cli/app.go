package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	appName      = "SSLB (github.com/eduardonunesp/sslb)"
	appUsage     = "sslb"
	versionMajor = "0"
	versionMinor = "1"
	versionBuild = "0"
)

func CreateAPP() {
	var rootCmd = &cobra.Command{
		Use: "sslb",
		Run: func(cmd *cobra.Command, args []string) {
			var verbose bool
			var filename string
			var err error

			fflags := cmd.Flags()
			verbose = fflags.Changed("verbose") == true

			if fflags.Changed("filename") {
				filename, err = fflags.GetString("filename")
				if err != nil {
					os.Exit(0)
					return
				}
			}

			RunServer(verbose, filename)
		},
	}

	rootCmd.Flags().BoolP("verbose", "v", false, "Help message for flag intone")
	rootCmd.Flags().StringP("filename", "f", "", "Set the filename as the configuration")

	statusCommand := &cobra.Command{
		Use: "status",
		Run: func(cmd *cobra.Command, args []string) {
			var filename string
			var err error

			fflags := cmd.Flags()

			if fflags.Changed("filename") {
				filename, err = fflags.GetString("filename")
				if err != nil {
					os.Exit(0)
					return
				}
			}

			InternalStatus(filename)
		},
	}

	statusCommand.Flags().StringP("filename", "f", "", "Set the filename as the configuration")

	rootCmd.AddCommand(statusCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
