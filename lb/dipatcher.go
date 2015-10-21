package lb

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Dispatcher struct {
	Mutex sync.Mutex
	Idle  bool
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{Idle: true}
}

func processReturn(result *http.Response) WorkRequest {
	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return NewWorkerRequest(http.StatusInternalServerError, result.Header, []byte(err.Error()))
	}

	return NewWorkerRequest(result.StatusCode, result.Header, []byte(body))
}

func execRequest(address string, r *http.Request) WorkRequest {
	var request *http.Request
	var err error

	client := &http.Client{}
	request, err = http.NewRequest(r.Method, address, r.Body)

	for k, vv := range r.Header {
		for _, v := range vv {
			request.Header.Set(k, v)
		}
	}

	response, err := client.Do(request)

	if err != nil {
		return NewWorkerRequestErr(http.StatusInternalServerError, []byte(err.Error()))
	}

	if response == nil {
		return NewWorkerRequestErr(http.StatusBadGateway, []byte("Method Not Supported By SSLB"))
	}

	return processReturn(response)
}

func (d *Dispatcher) Run(backend *Backend, r *http.Request, chanReceiver WorkRequestChan) {
	d.Mutex.Lock()
	d.Idle = false
	d.Mutex.Unlock()

	backend.Mutex.Lock()
	backend.Score += 1
	backend.Mutex.Unlock()

	chanReceiver <- execRequest(backend.Address, r)
	d.Mutex.Lock()
	d.Idle = true
	d.Mutex.Unlock()
}

type Dispatchers []*Dispatcher

type DispatcherPool struct {
	Mutex       sync.Mutex
	Size        int
	Dispatchers Dispatchers
}

func NewDispatcherPool(size int) *DispatcherPool {
	dp := &DispatcherPool{Size: size}
	dp.createPool()
	return dp
}

func (dp *DispatcherPool) createPool() {
	log.Printf("Create dispatcher pool with [%d]", dp.Size)
	for i := 0; i <= dp.Size; i++ {
		dispatcher := NewDispatcher()
		dp.Dispatchers = append(dp.Dispatchers, dispatcher)
	}
}

func (dp *DispatcherPool) Get(backend *Backend, r *http.Request, chanReceiver WorkRequestChan) {
	dp.Mutex.Lock()
	var idleDispatcher *Dispatcher

	for {
		for _, dispatcher := range dp.Dispatchers {
			dispatcher.Mutex.Lock()
			if dispatcher.Idle {
				dispatcher.Idle = false
				idleDispatcher = dispatcher
			}
			dispatcher.Mutex.Unlock()
		}

		if idleDispatcher != nil {
			break
		}

		time.Sleep(time.Millisecond)
	}

	idleDispatcher.Run(backend, r, chanReceiver)
	dp.Mutex.Unlock()
}
