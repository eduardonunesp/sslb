package worker

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/eduardonunesp/sslb/lb/endpoint"
	"github.com/eduardonunesp/sslb/lb/request"
)

type WorkerPool struct {
	Mutex   sync.Mutex
	Size    int
	Workers Workers
	DPPool  *DispatcherPool
}

func NewWorkerPool(wSize, dSize int) *WorkerPool {
	wp := &WorkerPool{Size: wSize}
	wp.DPPool = NewDispatcherPool(dSize)
	wp.createPool()
	return wp
}

func (wp *WorkerPool) createPool() {
	log.Printf("Create worker pool with [%d]", wp.Size)
	for i := 0; i <= wp.Size; i++ {
		worker := NewWorker(wp.DPPool)
		wp.Workers = append(wp.Workers, worker)
	}
}

func (wp *WorkerPool) Get(r *http.Request, frontend *endpoint.Frontend) request.SSLBRequestChan {
	wp.Mutex.Lock()
	var idleWorker *Worker

	for {

		for _, worker := range wp.Workers {
			worker.Mutex.Lock()
			if worker.Idle {
				worker.Idle = false
				idleWorker = worker
				worker.Mutex.Unlock()
				break
			}
			worker.Mutex.Unlock()
		}

		if idleWorker != nil {
			break
		}

		time.Sleep(time.Millisecond)
	}

	c := idleWorker.Run(r, frontend)
	wp.Mutex.Unlock()
	return c
}

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

func (dp *DispatcherPool) Get(backend *endpoint.Backend, r *http.Request, chanReceiver request.SSLBRequestChan) {
	dp.Mutex.Lock()
	var idleDispatcher *Dispatcher

	for {
		for _, dispatcher := range dp.Dispatchers {
			dispatcher.Mutex.Lock()
			if dispatcher.Idle {
				dispatcher.Idle = false
				idleDispatcher = dispatcher
				dispatcher.Mutex.Unlock()
				break
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
