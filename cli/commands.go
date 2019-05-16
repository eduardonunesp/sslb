package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/eduardonunesp/sslb/cfg"
	"github.com/eduardonunesp/sslb/lb"
	sslbRPC "github.com/eduardonunesp/sslb/rpc"
	"github.com/olekukonko/tablewriter"
)

const (
	CONFIG_FILENAME         = "config.json"
	CONFIG_FILENAME_EXAMPLE = "config.json.example"
)

func InternalStatus(filename string) {
	if filename == "" {
		filename = CONFIG_FILENAME
	}

	configuration := cfg.Setup(filename)
	address := fmt.Sprintf("%s:%d",
		configuration.GeneralConfig.RPCHost,
		configuration.GeneralConfig.RPCPort,
	)

	log.Println("Start SSLB (Client) ")

	client, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	reply := sslbRPC.StatusResponse{}

	rpcCall := jsonrpc.NewClient(client)
	err = rpcCall.Call("ServerStatus.GetIdle", 0, &reply)

	if err != nil {
		log.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Workers Idle"})
	idles := fmt.Sprintf("%d", reply.IdleWPool)
	table.Append([]string{idles})
	table.Render()
}

func RunServer(verbose bool, filename string) {
	if !verbose {
		log.SetOutput(ioutil.Discard)
	}

	if filename == "" {
		filename = CONFIG_FILENAME
	}

	log.Println("Start SSLB (Server) ")
	log.Println("Using config:", filename)

	// The function setup do everything for configure
	// and return the server ready to run
	configuration := cfg.Setup(filename)
	server := lb.NewServer(configuration)
	sslbRPC.StartServer(server)

	log.Println("Prepare to run server ...")
	server.Run()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	server.Stop()
}
