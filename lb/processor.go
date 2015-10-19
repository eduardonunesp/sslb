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
	errTimeout        = errors.New("Timeout")
	errInvalidRequest = errors.New("Invalid request")
)

type Processor struct {
	Frontends Frontends
}

func NewProcessor() *Processor {
	runtime.GOMAXPROCS(1)
	return &Processor{}
}

func (p *Processor) Status() {
	go func() {
		for {
			log.Println("--------- Status ---------")
			for _, frontend := range p.Frontends {
				log.Println("%T", frontend)
				for _, backend := range frontend.Backends {
					log.Println("%T", backend)
				}
			}

			time.Sleep(time.Second * 5)
		}
	}()
}

func (p *Processor) checkIfExists(route string) bool {
	exists := false
	for _, frontend := range p.Frontends {
		if frontend.Route == route {
			exists = true
		}
	}

	return exists
}

func (p *Processor) AddFrontend(frontend *Frontend) {
	if !p.checkIfExists(frontend.Route) {
		log.Println("Route added", frontend.Route)
		p.Frontends = append(p.Frontends, frontend)
	} else {
		log.Println("Route already exists for", frontend.Route)
	}
}

func (p *Processor) RunFrontendProcessor(frontend *Frontend) {
	host := frontend.Host
	port := frontend.Port
	address := fmt.Sprintf("%s:%d", host, port)

	log.Println("Run frontend processor at", address)

	http.HandleFunc(frontend.Route, func(w http.ResponseWriter, r *http.Request) {
		chanResponse := WorkerRun(r, frontend)
		defer close(chanResponse)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		select {
		case result := <-chanResponse:
			if result.Status > 400 {
				http.Error(w, string(result.Body), result.Status)
			} else {
				w.WriteHeader(result.Status)
				w.Write(result.Body)
			}

		case <-ticker.C:
			http.Error(w, errTimeout.Error(), http.StatusRequestTimeout)
		}
	})

	http.ListenAndServe(address, nil)
}

func (p *Processor) Run(processor *Processor) {
	p.Status()
	for _, frontend := range p.Frontends {
		p.RunFrontendProcessor(frontend)
	}
}
