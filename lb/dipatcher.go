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

func DispatchRequest(address string, r *http.Request, chanReceiver WorkRequestChan) {
	if r.Method == "GET" {
		chanReceiver <- getRequest(address)
	}
}
