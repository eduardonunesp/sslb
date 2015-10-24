package worker

import (
	"io/ioutil"
	"net/http"
	"strings"
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
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return request.NewWorkerRequest(http.StatusInternalServerError, result.Header, []byte(err.Error()))
	}

	return request.NewWorkerRequest(result.StatusCode, result.Header, []byte(body))
}

func checkForWebsocket(r *http.Request) bool {
	result := false
	connHDR := ""
	connHDRS := r.Header["Connection"]

	if len(connHDRS) > 0 {
		connHDR = connHDRS[0]
	}

	if connHDR == "upgrade" || connHDR == "Upgrade" {
		upgradeHDRS := r.Header["Upgrade"]
		if len(upgradeHDRS) > 0 {
			result = (strings.ToLower(upgradeHDRS[0]) == "websocket")
		}
	}

	return result
}

func execRequest(address string, r *http.Request) request.SSLBRequest {
	var httpRequest *http.Request
	var err error

	if checkForWebsocket(r) {
		ret := request.NewWorkerRequestUpgraded()
		ret.Address = address
		return ret
	}

	requestAddress := address + r.URL.String()

	client := &http.Client{}
	httpRequest, err = http.NewRequest(r.Method, requestAddress, r.Body)

	for k, vv := range r.Header {
		for _, v := range vv {
			httpRequest.Header.Set(k, v)
		}
	}

	response, err := client.Do(httpRequest)

	if err != nil {
		return request.NewWorkerRequestErr(http.StatusRequestTimeout, []byte("No backend available"))
	}

	if response == nil {
		return request.NewWorkerRequestErr(http.StatusBadGateway, []byte("Method Not Supported By SSLB"))
	}

	ret := processReturn(response)
	ret.Address = address
	return ret
}

func (d *Dispatcher) Run(backend *endpoint.Backend, r *http.Request, chanReceiver request.SSLBRequestChan) {
	d.Mutex.Lock()
	d.Idle = false
	d.Mutex.Unlock()

	backend.Mutex.Lock()
	backend.Score++
	backend.Mutex.Unlock()

	go func(c request.SSLBRequestChan) {
		// On a serious problem
		defer func() {
			if rec := recover(); rec != nil {
				// Channel not used
			}
		}()

		c <- execRequest(backend.Address, r)
	}(chanReceiver)

	d.Mutex.Lock()
	d.Idle = true
	d.Mutex.Unlock()
}
