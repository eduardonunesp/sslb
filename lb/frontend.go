package lb

import "time"

const (
	StrategyEqual = iota
	StrategyRoundRobin
)

type Frontend struct {
	Host  string
	Port  uint
	Route string

	Backends Backends
	Strategy uint
	Timeout  time.Duration
}

type Frontends []*Frontend

func NewFrontend(host string, port uint, route string) *Frontend {
	return &Frontend{
		Host:     host,
		Port:     port,
		Route:    route,
		Strategy: StrategyEqual,
		Timeout:  time.Millisecond * 1000 * 5,
	}
}

func (f *Frontend) AddBackend(backend *Backend) {
	f.Backends = append(f.Backends, backend)
}

func (f *Frontend) SetStrategy(balance uint) {
	f.Strategy = balance
}
