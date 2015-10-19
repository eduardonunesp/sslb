package lb

const (
	BalanceEqual = iota
)

type Frontend struct {
	Host     string
	Port     uint
	Route    string
	Backends Backends
	Balance  uint
}

type Frontends []*Frontend

func NewFrontend(host string, port uint, route string) *Frontend {
	return &Frontend{Host: host, Port: port, Route: route}
}

func (f *Frontend) AddBackend(backend *Backend) {
	f.Backends = append(f.Backends, backend)
}

func (f *Frontend) TestRoute(route string) bool {
	if route == f.Route {
		return true
	} else {
		return false
	}
}

func (f *Frontend) SetBalance(balance uint) {
	f.Balance = balance
}
