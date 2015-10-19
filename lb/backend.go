package lb

import "sync"

type Backend struct {
	Address string
	Active  bool
	Tries   uint
	Score   uint
	Mutex   sync.Mutex
}

type Backends []*Backend

type ByScore []*Backend

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].Score < a[j].Score }

func NewBackend(address string) *Backend {
	return &Backend{
		Address: address,
		Active:  true,
		Tries:   0,
		Score:   0,
	}
}
