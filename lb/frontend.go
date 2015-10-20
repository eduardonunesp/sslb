package lb

import (
	"errors"
	"time"
)

var (
	errInvalidFrontend = errors.New("Invalid frontend configuration")
)

const (
	StrategyEqual = iota
	StrategyRoundRobin
)

type Frontend struct {
	Name  string
	Host  string
	Port  uint
	Route string

	Backends Backends
	Strategy uint
	Timeout  time.Duration
	WPool    *WorkerPool
}

type Frontends []*Frontend

func NewFrontend(name string, host string, port uint, route string) *Frontend {
	wp := NewWorkerPool(100, 100)
	return &Frontend{
		Name:     name,
		Host:     host,
		Port:     port,
		Route:    route,
		Strategy: StrategyEqual,
		Timeout:  time.Millisecond * 1000 * 5,
		WPool:    wp,
	}
}

func (f *Frontend) AddBackend(backend *Backend) {
	f.Backends = append(f.Backends, backend)
}

func (f *Frontend) SetStrategy(balance uint) {
	f.Strategy = balance
}
