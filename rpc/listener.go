package rpc

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/eduardonunesp/sslb/lb"
)

type ServerStatus struct {
	Server *lb.Server
}

type StatusResponse struct {
	IdleWPool int
}

func (s *ServerStatus) GetIdle(args interface{}, reply *StatusResponse) error {
	statusRes := StatusResponse{
		s.Server.CountIdle(),
	}

	*reply = statusRes
	return nil
}

func StartServer(s *lb.Server) {
	go func() {
		serverStatus := &ServerStatus{s}

		server := rpc.NewServer()
		server.Register(serverStatus)

		address := fmt.Sprintf("%s:%d",
			s.Configuration.GeneralConfig.RPCHost,
			s.Configuration.GeneralConfig.RPCPort,
		)

		server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
		listener, e := net.Listen("tcp", address)
		if e != nil {
			log.Fatal("listen error:", e)
		}

		for {
			if conn, err := listener.Accept(); err != nil {
				log.Fatal("accept error: " + err.Error())
			} else {
				go server.ServeCodec(jsonrpc.NewServerCodec(conn))
			}
		}
	}()
}
