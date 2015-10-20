package main

import "github.com/eduardonunesp/sslb/lb"

func main() {
	server := lb.NewServer(4)

	frontend := lb.NewFrontend("front1", "localhost", 9000, "/")
	frontend.AddBackend(lb.NewBackend("back1", "http://localhost:9001"))
	frontend.AddBackend(lb.NewBackend("back2", "http://localhost:9002"))

	server.AddFrontend(frontend)

	server.Run()
}
