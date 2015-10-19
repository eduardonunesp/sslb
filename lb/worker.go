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

func preProcessWorker(frontend *Frontend) *Backend {
	backendsSlice := []*Backend{}

	for _, backend := range frontend.Backends {
		backendsSlice = append(backendsSlice, backend)
	}

	sort.Sort(ByScore(backendsSlice))

	backend := backendsSlice[0]
	return backend
}

// Run the worker to request a job
func WorkerRun(r *http.Request, frontend *Frontend) WorkRequestChan {
	chanReceiver := make(WorkRequestChan)
	go func() {
		backend := preProcessWorker(frontend)
		DispatchRequest(backend, r, chanReceiver)
	}()
	return chanReceiver
}
