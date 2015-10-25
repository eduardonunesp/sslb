package worker

import (
	"net/http"
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
	frontend.Mutex.Lock()
	defer frontend.Mutex.Unlock()

	var backendWithMinScore *endpoint.Backend

	for idx, backend := range frontend.Backends {
		backend.RWMutex.RLock()
		if idx == 0 {
			backendWithMinScore = backend
		} else {
			if backend.Score < backendWithMinScore.Score {
				backendWithMinScore = backend
			}
		}
		backend.RWMutex.RUnlock()
	}

	return backendWithMinScore
}

func (w *Worker) Run(r *http.Request, frontend *endpoint.Frontend) request.SSLBRequestChan {
	w.Mutex.Lock()
	w.Idle = false
	w.Mutex.Unlock()

	chanReceiver := make(request.SSLBRequestChan)
	go func(w *Worker, chanReceiver request.SSLBRequestChan, f *endpoint.Frontend) {
		backend := preProcessWorker(f)

		if backend != nil {
			backend.RWMutex.Lock()
			backend.Score++
			backend.RWMutex.Unlock()

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
