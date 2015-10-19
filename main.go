package main

import "github.com/eduardonunesp/sslb/lb"

func main() {
	processor := lb.NewProcessor()

	frontend := lb.NewFrontend("localhost", 9000, "/")
	frontend.AddBackend(lb.NewBackend("back1", "http://localhost:9001"))
	frontend.AddBackend(lb.NewBackend("back2", "http://localhost:9002"))

	processor.AddFrontend(frontend)

	processor.Run(processor)
}
