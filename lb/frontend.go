package lb

import (
	"errors"
	"time"
)

var (
	errInvalidFrontend = errors.New("Invalid frontend configuration")
)

const (
	StrategyRoundRobin = iota
)

// Frontend structure
type Frontend struct {
	Name  string
	Host  string
	Port  int
	Route string

	Backends Backends
	Strategy int
	Timeout  time.Duration
	WPool    *WorkerPool
}

type Frontends []*Frontend

func NewFrontend(name string, host string,
	port int, route string, timeout int,
	workerPoolSize, dispatcherPoolSize int) *Frontend {

	// Config the pool size for workers and dispatchers
	wp := NewWorkerPool(workerPoolSize, dispatcherPoolSize)
	return &Frontend{
		Name:    name,
		Host:    host,
		Port:    port,
		Route:   route,
		Timeout: time.Duration(timeout) * time.Millisecond,

		Strategy: StrategyRoundRobin,
		WPool:    wp,
	}
}

func (f *Frontend) AddBackend(backend *Backend) {
	f.Backends = append(f.Backends, backend)
}

func (f *Frontend) SetStrategy(balance int) {
	f.Strategy = balance
}
