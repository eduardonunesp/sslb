package lb

import (
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

type WorkRequest struct {
	Status   int
	Header   http.Header
	Body     []byte
	Internal bool
}

type WorkRequestChan chan WorkRequest

func NewWorkerRequestErr(status int, body []byte) WorkRequest {
	return WorkRequest{Status: status, Body: body, Internal: true}
}

func NewWorkerRequest(status int, header http.Header, body []byte) WorkRequest {
	return WorkRequest{Status: status, Header: header, Body: body}
}

type WorkerFunc func(*http.Request, *Frontend) WorkRequestChan

type Worker struct {
	Mutex sync.Mutex
	Idle  bool
	DPool *DispatcherPool
}

func NewWorker(dp *DispatcherPool) *Worker {
	return &Worker{Idle: true, DPool: dp}
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

func (w *Worker) Run(r *http.Request, frontend *Frontend) WorkRequestChan {
	w.Mutex.Lock()
	w.Idle = false
	w.Mutex.Unlock()

	chanReceiver := make(WorkRequestChan)
	go func(w *Worker, chanReceiver WorkRequestChan, frontend *Frontend) {
		backend := preProcessWorker(frontend)

		if backend != nil {
			w.DPool.Get(backend, r, chanReceiver)
		} else {
			chanReceiver <- NewWorkerRequestErr(500, []byte("No backend available"))
		}

		w.Mutex.Lock()
		w.Idle = true
		w.Mutex.Unlock()
	}(w, chanReceiver, frontend)

	return chanReceiver
}

type Workers []*Worker

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

func (wp *WorkerPool) Get(r *http.Request, frontend *Frontend) WorkRequestChan {
	wp.Mutex.Lock()
	var idleWorker *Worker

	for {

		for _, worker := range wp.Workers {
			worker.Mutex.Lock()
			if worker.Idle {
				worker.Idle = false
				idleWorker = worker
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
