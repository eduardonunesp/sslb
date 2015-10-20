package lb

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

var (
	errNoFrontend     = errors.New("No frontend configuration detected")
	errNoBackend      = errors.New("No backend configuration detected")
	errTimeout        = errors.New("Timeout")
	errInvalidRequest = errors.New("Invalid request")
	errPortExists     = errors.New("Port already in use")
	errRouteExists    = errors.New("Route already in use")
)

type Server struct {
	Frontends Frontends
}

func NewServer(procs int) *Server {
	runtime.GOMAXPROCS(procs)
	return &Server{}
}

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

func (s *Server) AddFrontend(frontend *Frontend) {
	err := s.preChecksBeforeAdd(frontend)
	if err != nil {
		log.Fatal(err.Error())
	}

	s.Frontends = append(s.Frontends, frontend)
}

func (s *Server) RunFrontendServer(frontend *Frontend) {
	if len(frontend.Backends) == 0 {
		log.Fatal(errNoBackend.Error())
	}

	host := frontend.Host
	port := frontend.Port
	address := fmt.Sprintf("%s:%d", host, port)

	for _, backend := range frontend.Backends {
		backend.HeartCheck()
	}

	log.Printf("Run frontend server [%s] at [%s]", frontend.Name, address)

	httpHandle := http.NewServeMux()

	httpHandle.HandleFunc(frontend.Route, func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Println("Err", rec)
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()

		chanResponse := frontend.WPool.Get(r, frontend)
		defer close(chanResponse)

		r.Close = true

		ticker := time.NewTicker(frontend.Timeout)
		defer ticker.Stop()

		select {
		case result := <-chanResponse:
			if result.Status >= 400 {
				http.Error(w, string(result.Body), result.Status)
			} else {
				w.WriteHeader(result.Status)
				w.Write(result.Body)
			}

		case <-ticker.C:
			http.Error(w, errTimeout.Error(), http.StatusRequestTimeout)
		}
	})

	server := &http.Server{
		Addr:    address,
		Handler: httpHandle,
	}

	server.ListenAndServe()
}

func (s *Server) Run() {
	if len(s.Frontends) == 0 {
		log.Fatal(errNoFrontend.Error())
	}

	for _, frontend := range s.Frontends {
		s.RunFrontendServer(frontend)
	}
}
