package lb

import (
	"net/http"
	"sort"
)

type WorkRequest struct {
	Status int
	Body   []byte
}

type WorkRequestChan chan WorkRequest

func NewWorkerRequest(status int, result []byte) WorkRequest {
	return WorkRequest{status, result}
}

// Search for backend with the less score
func preProcessWorker(frontend *Frontend) *Backend {
	backendsSlice := []*Backend{}

	for _, backend := range frontend.Backends {
		if backend.Active && !backend.Failed {
			backendsSlice = append(backendsSlice, backend)
		}
	}

	sort.Sort(ByScore(backendsSlice))

	var backend *Backend
	if len(backendsSlice) > 0 {
		backend = backendsSlice[0]
	}

	return backend
}

// Run the worker to request a job
func WorkerRun(r *http.Request, frontend *Frontend) WorkRequestChan {
	chanReceiver := make(WorkRequestChan)
	go func() {
		backend := preProcessWorker(frontend)

		if backend != nil {
			DispatchRequest(backend, r, chanReceiver)
		} else {
			chanReceiver <- NewWorkerRequest(500, []byte("No backend available"))
		}
	}()
	return chanReceiver
}
