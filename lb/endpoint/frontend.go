package endpoint

import (
	"sync"
	"time"
)

const (
	// StrategyRoundRobin is the strategy available for now
	StrategyRoundRobin = iota
)

// Frontend structure
type Frontend struct {
	Mutex sync.Mutex

	Name  string
	Host  string
	Port  int
	Route string

	Backends Backends
	Strategy int
	Timeout  time.Duration
}

type Frontends []*Frontend

// Create and returns a new Frontend
func NewFrontend(name string, host string,
	port int, route string, timeout int) *Frontend {

	return &Frontend{
		Name:    name,
		Host:    host,
		Port:    port,
		Route:   route,
		Timeout: time.Duration(timeout) * time.Millisecond,

		Strategy: StrategyRoundRobin,
	}
}

// AddBackend it's responsible to link a backend conf with frontend
func (f *Frontend) AddBackend(backend *Backend) {
	f.Backends = append(f.Backends, backend)
}

// SetStrategy will set the strategy of balancing
func (f *Frontend) SetStrategy(balance int) {
	f.Strategy = balance
}
