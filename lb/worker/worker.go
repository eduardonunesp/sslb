package worker

import (
	"net/http"
	"sort"
	"sync"

	"github.com/eduardonunesp/sslb/lb/endpoint"
	"github.com/eduardonunesp/sslb/lb/request"
)

type Worker struct {
	Mutex sync.Mutex
	Idle  bool
	DPool *DispatcherPool
}

type Workers []*Worker

func NewWorker(dp *DispatcherPool) *Worker {
	return &Worker{Idle: true, DPool: dp}
}

// Search for backend with the less score
func preProcessWorker(frontend *endpoint.Frontend) *endpoint.Backend {
	backendsSlice := []*endpoint.Backend{}

	for _, backend := range frontend.Backends {
		backend.Mutex.Lock()
		if backend.Active && !backend.Failed {
			backendsSlice = append(backendsSlice, backend)
		}
		backend.Mutex.Unlock()
	}

	sort.Sort(endpoint.ByScore(backendsSlice))

	var backend *endpoint.Backend
	if len(backendsSlice) > 0 {
		backend = backendsSlice[0]
	}

	return backend
}

func (w *Worker) Run(r *http.Request, frontend *endpoint.Frontend) request.SSLBRequestChan {
	w.Mutex.Lock()
	w.Idle = false
	w.Mutex.Unlock()

	chanReceiver := make(request.SSLBRequestChan)
	go func(w *Worker, chanReceiver request.SSLBRequestChan, frontend *endpoint.Frontend) {
		backend := preProcessWorker(frontend)

		if backend != nil {
			w.DPool.Get(backend, r, chanReceiver)
		} else {
			chanReceiver <- request.NewWorkerRequestErr(http.StatusServiceUnavailable, []byte("Service Unavailable"))
		}

		w.Mutex.Lock()
		w.Idle = true
		w.Mutex.Unlock()
	}(w, chanReceiver, frontend)

	return chanReceiver
}
