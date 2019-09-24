package impl

import (
	"errors"
	"log"
	"runtime"

	"github.com/eduardonunesp/sslb/types"
)

var (
	errNoFrontend  = errors.New("No frontend configuration detected")
	errNoBackend   = errors.New("No backend configuration detected")
	errTimeout     = errors.New("Timeout")
	errPortExists  = errors.New("Port already in use")
	errRouteExists = errors.New("Route already in use")
)

type Server struct {
	config                 types.Config
	frontendManagerFactory types.FrontendManagerFactoryCreator
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SetConfig(config types.Config) {
	s.config = config
}

func (s *Server) SetFrontendManagerFactory(frontendManagerFactory types.FrontendManagerFactoryCreator) {
	s.frontendManagerFactory = frontendManagerFactory
}

func (s *Server) Run() {
	frontendManager := s.frontendManagerFactory.CreateNewFrontendManager()
	runtime.GOMAXPROCS(s.config.GeneralConfig.MaxProcs)
	frontendManager.SetConfig(s.config.FrontendConfigList)
	frontendManager.LoadConfig()
	frontendManager.Run()
}

func (s *Server) Stop() {
	log.Println("Bye")
}
