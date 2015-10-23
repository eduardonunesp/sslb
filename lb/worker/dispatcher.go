package worker

import (
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/eduardonunesp/sslb/lb/endpoint"
	"github.com/eduardonunesp/sslb/lb/request"
)

type Dispatcher struct {
	Mutex sync.Mutex
	Idle  bool
}

type Dispatchers []*Dispatcher

func NewDispatcher() *Dispatcher {
	return &Dispatcher{Idle: true}
}

func processReturn(result *http.Response) request.SSLBRequest {
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return request.NewWorkerRequest(http.StatusInternalServerError, result.Header, []byte(err.Error()))
	}

	return request.NewWorkerRequest(result.StatusCode, result.Header, []byte(body))
}

func execRequest(address string, r *http.Request) request.SSLBRequest {
	var httpRequest *http.Request
	var err error

	requestAddress := address + r.URL.String()

	client := &http.Client{}
	httpRequest, err = http.NewRequest(r.Method, requestAddress, r.Body)

	for k, vv := range r.Header {
		for _, v := range vv {
			httpRequest.Header.Set(k, v)
		}
	}

	response, err := client.Do(httpRequest)
	defer response.Body.Close()

	if err != nil {
		return request.NewWorkerRequestErr(http.StatusRequestTimeout, []byte("No backend available"))
	}

	if response == nil {
		return request.NewWorkerRequestErr(http.StatusBadGateway, []byte("Method Not Supported By SSLB"))
	}

	return processReturn(response)
}

func (d *Dispatcher) Run(backend *endpoint.Backend, r *http.Request, chanReceiver request.SSLBRequestChan) {
	d.Mutex.Lock()
	d.Idle = false
	d.Mutex.Unlock()

	backend.Mutex.Lock()
	backend.Score++
	backend.Mutex.Unlock()

	chanReceiver <- execRequest(backend.Address, r)
	d.Mutex.Lock()
	d.Idle = true
	d.Mutex.Unlock()
}
