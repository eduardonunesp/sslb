package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/eduardonunesp/sslb/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/eduardonunesp/sslb/cfg"
	sslbRPC "github.com/eduardonunesp/sslb/rpc"
)

func InternalStatus(c *cli.Context) {
	client, err := net.Dial("tcp", "127.0.0.1:45222")
	if err != nil {
		log.Fatal(err)
	}

	reply := sslbRPC.StatusResponse{}

	rpcCall := jsonrpc.NewClient(client)
	err = rpcCall.Call("ServerStatus.GetIdle", 0, &reply)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result: %d %d\n", reply.IdleWPool, reply.IdleDPool)
}

func RunServer(c *cli.Context) {
	if !c.Bool("verbose") {
		log.SetOutput(ioutil.Discard)
	}

	if c.Bool("config") {
		cfg.CreateConfig(CONFIG_FILENAME)
		os.Exit(0)
	}

	filename := "config.json"
	if c.String("filename") != "" {
		filename = c.String("filename")
	}

	// The function setup do everything for configure
	// and return the server ready to run
	server := cfg.Setup(filename)
	sslbRPC.StartServer(server)

	server.Run()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	server.Stop()
}
