package lb

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/eduardonunesp/sslb/lb/endpoint"
	"github.com/eduardonunesp/sslb/lb/worker"
)

var (
	errNoFrontend  = errors.New("No frontend configuration detected")
	errNoBackend   = errors.New("No backend configuration detected")
	errTimeout     = errors.New("Timeout")
	errPortExists  = errors.New("Port already in use")
	errRouteExists = errors.New("Route already in use")
)

type Server struct {
	Ch        chan bool
	Frontends endpoint.Frontends
	WPool     *worker.WorkerPool
	WaitGroup *sync.WaitGroup
}

func NewServer(workerPoolSize, dispatcherPoolSize int) *Server {
	// Config the pool size for workers and dispatchers
	wp := worker.NewWorkerPool(workerPoolSize, dispatcherPoolSize)
	s := &Server{
		Ch:        make(chan bool),
		WPool:     wp,
		WaitGroup: &sync.WaitGroup{},
	}

	return s
}

// Some previous checkings before run
func (s *Server) preChecksBeforeAdd(newFrontend *endpoint.Frontend) error {
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

func (s *Server) AddFrontend(frontend *endpoint.Frontend) {
	err := s.preChecksBeforeAdd(frontend)
	if err != nil {
		log.Fatal(err.Error())
	}

	s.Frontends = append(s.Frontends, frontend)
}

// Lets run the frontned
func (s *Server) RunFrontendServer(frontend *endpoint.Frontend) {
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
		s.WaitGroup.Add(1)
		defer s.WaitGroup.Done()

		// On a serious problem
		defer func() {
			if rec := recover(); rec != nil {
				log.Println("Err", rec)
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()

		// Get a channel the already attached to a worker
		chanResponse := s.WPool.Get(r, frontend)
		defer close(chanResponse)

		r.Close = true

		// Timeout ticker
		ticker := time.NewTicker(frontend.Timeout)
		defer ticker.Stop()

		select {
		case result := <-chanResponse:
			// We have a response, it's valid ?
			if result.Internal {
				http.Error(w, string(result.Body), result.Status)
			} else {
				for k, vv := range result.Header {
					for _, v := range vv {
						w.Header().Set(k, v)
					}
				}

				w.WriteHeader(result.Status)
				w.Write(result.Body)
			}
		case <-r.Cancel:
			// Request cancelled
		// Timeout
		case <-ticker.C:
			http.Error(w, errTimeout.Error(), http.StatusRequestTimeout)
		}
	})

	// Config and start server
	server := &http.Server{
		Addr:    address,
		Handler: httpHandle,
	}

	server.ListenAndServe()
}

func (s *Server) Stop() {
	close(s.Ch)
	s.WaitGroup.Wait()
}

func (s *Server) Run() {
	if len(s.Frontends) == 0 {
		log.Fatal(errNoFrontend.Error())
	}

	// Run the fronend config
	for _, frontend := range s.Frontends {
		go s.RunFrontendServer(frontend)
	}
}
