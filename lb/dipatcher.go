package lb

import (
	"io/ioutil"
	"net/http"
)

func getRequest(address string) WorkRequest {
	result, err := http.Get(address)
	if err != nil {
		return NewWorkerRequest(http.StatusInternalServerError, []byte("Address is out of reach"))
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return NewWorkerRequest(http.StatusInternalServerError, []byte(err.Error()))
	}

	return NewWorkerRequest(http.StatusOK, []byte(body))
}

func DispatchRequest(backend *Backend, r *http.Request, chanReceiver WorkRequestChan) {
	if r.Method == "GET" {
		backend.Mutex.Lock()
		backend.Score += 1
		backend.Mutex.Unlock()
		chanReceiver <- getRequest(backend.Address)
	}
}
