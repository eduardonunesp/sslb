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

func getRequest(address string) WorkRequest {
	result, err := http.Get(address)
	if err != nil {
		return NewWorkerRequest(http.StatusInternalServerError, []byte(err.Error()))
	}

	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return NewWorkerRequest(http.StatusInternalServerError, []byte(err.Error()))
	}

	return NewWorkerRequest(http.StatusOK, []byte(body))
}

func (d *Dispatcher) Run(backend *Backend, r *http.Request, chanReceiver WorkRequestChan) {
	d.Mutex.Lock()
	d.Idle = false
	d.Mutex.Unlock()

	if r.Method == "GET" {
		backend.Mutex.Lock()
		backend.Score += 1
		backend.Mutex.Unlock()

		go func(c WorkRequestChan, d *Dispatcher) {
			c <- getRequest(backend.Address)
			d.Idle = true
		}(chanReceiver, d)
	}
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
			if dispatcher.Idle {
				dispatcher.Mutex.Lock()
				dispatcher.Idle = false
				idleDispatcher = dispatcher
				dispatcher.Mutex.Unlock()
				break
			}
		}

		if idleDispatcher != nil {
			break
		}

		time.Sleep(time.Millisecond)
	}

	dp.Mutex.Unlock()

	idleDispatcher.Run(backend, r, chanReceiver)
}
