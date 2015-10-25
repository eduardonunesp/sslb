package rpc

import (
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
	IdleDPool int
}

func (s *ServerStatus) GetIdle(args interface{}, reply *StatusResponse) error {
	statusRes := StatusResponse{
		s.Server.WPool.CountIdle(),
		s.Server.WPool.DPPool.CountIdle(),
	}

	*reply = statusRes
	return nil
}

func StartServer(s *lb.Server) {
	go func() {
		serverStatus := &ServerStatus{s}

		server := rpc.NewServer()
		server.Register(serverStatus)

		server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
		listener, e := net.Listen("tcp", "127.0.0.1:45222")
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
