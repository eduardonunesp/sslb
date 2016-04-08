package lb

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var (
	errNoFrontend  = errors.New("No frontend configuration detected")
	errNoBackend   = errors.New("No backend configuration detected")
	errTimeout     = errors.New("Timeout")
	errPortExists  = errors.New("Port already in use")
	errRouteExists = errors.New("Route already in use")
)

type ShutdownChan chan bool

type Server struct {
	Configuration
	Frontends
	ShutdownChan
	*WorkerPool

	sync.Mutex
	*sync.WaitGroup
}

func NewServer(configuration Configuration) *Server {
	return &Server{
		Configuration: configuration,
		ShutdownChan:  make(ShutdownChan),
		WaitGroup:     &sync.WaitGroup{},
		WorkerPool:    NewWorkerPool(configuration),
	}
}

func (s *Server) setup() {
	runtime.GOMAXPROCS(s.Configuration.GeneralConfig.MaxProcs)

	for _, frontend := range s.Configuration.FrontendsConfig {
		_ = frontend
		newFrontend := NewFrontend(frontend)
		for _, backend := range frontend.BackendsConfig {
			newFrontend.Backends = append(newFrontend.Backends, NewBackend(backend))
		}

		if err := s.preChecksBeforeAdd(newFrontend); err != nil {
			log.Fatal(err.Error())
		} else {
			s.Frontends = append(s.Frontends, newFrontend)
		}
	}
}

// Some previous checkings before run
func (s *Server) preChecksBeforeAdd(newFrontend *Frontend) error {
	for _, frontend := range s.Frontends {
		if frontend.Route == newFrontend.Route {
			return errRouteExists
		}

		if frontend.Port == newFrontend.Port {
			return errPortExists
		}

		if len(newFrontend.Backends) == 0 {
			return errNoBackend
		}
	}

	return nil
}

// Lets run the frontend
func (s *Server) RunFrontendServer(frontend *Frontend) {
	if len(frontend.Backends) == 0 {
		log.Fatal(errNoBackend.Error())
	}

	host := frontend.Host
	port := frontend.Port
	address := fmt.Sprintf("%s:%d", host, port)

	for _, backend := range frontend.Backends {
		// Before start the backend let's set a monitor
		backend.HeartCheck()
	}

	log.Printf("Run frontend server [%s] at [%s]", frontend.Name, address)

	// Prepare the mux
	httpHandle := http.NewServeMux()

	httpHandle.HandleFunc(frontend.Route, func(w http.ResponseWriter, r *http.Request) {
		s.Lock()
		s.Add(1)
		s.Unlock()

		// On a serious problem
		defer func() {
			if rec := recover(); rec != nil {
				log.Println("Err", rec)
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()

		bResponse := s.Get(r, frontend)
		r.Close = true

		// Timeout ticker
		ticker := time.NewTicker(frontend.Timeout)
		defer ticker.Stop()

		// We have a response, it's valid ?
		for k, vv := range bResponse.Header {
			for _, v := range vv {
				w.Header().Set(k, v)
			}
		}

		if bResponse.Upgraded {
			if s.Configuration.GeneralConfig.Websocket {
				bResponse.HijackWebSocket(w, r)
			}
		} else {
			w.WriteHeader(bResponse.Status)

			if r.Method != "HEAD" {
				w.Write(bResponse.Body)
			}
		}

		s.Lock()
		s.Done()
		s.Unlock()
	})

	// Config and start server
	server := &http.Server{
		Addr:    address,
		Handler: httpHandle,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Run() {
	log.Println("Setup and check configuration")
	s.setup()

	if len(s.Frontends) == 0 {
		log.Fatal(errNoFrontend.Error())
	}

	log.Println("Setup ok ...")

	// Run the fronend config
	for _, frontend := range s.Frontends {
		go s.RunFrontendServer(frontend)
	}
}

func (s *Server) Stop() {
	if s.Configuration.GeneralConfig.GracefulShutdown {
		log.Println("Wait for graceful shutdown")
		s.Wait()
		log.Println("Bye")
	}

	close(s.ShutdownChan)
}
