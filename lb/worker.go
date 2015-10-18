package lb

import "net/http"

type WorkRequest struct {
	Status int
	Body   []byte
}

type WorkRequestChan chan WorkRequest

func NewWorkerRequest(status int, result []byte) WorkRequest {
	return WorkRequest{status, result}
}

// Run the worker to request a job
func WorkerRun(r *http.Request) WorkRequestChan {
	chanReceiver := make(WorkRequestChan)
	go func() {
		DispatchRequest("http://localhost:9001", r, chanReceiver)
	}()
	return chanReceiver
}
