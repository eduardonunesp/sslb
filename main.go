package main

import "github.com/eduardonunesp/sslb/cfg"

func main() {
	// The function setup do everything for configure
	// and return the server ready to run
	server := cfg.Setup()
	server.Run()
}
