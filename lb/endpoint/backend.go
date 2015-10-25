package endpoint

import (
	"log"
	"net/http"
	"sync"
	"time"
)

//TODO: Need to rebalance the score when backend back to active

// Backend structure
type Backend struct {
	RWMutex sync.RWMutex

	Name      string
	Address   string
	Heartbeat string
	HBMethod  string

	ActiveAfter   int
	InactiveAfter int
	// Consider inactive after max inactiveAfter

	HeartbeatTime time.Duration // Heartbeat time if health
	RetryTime     time.Duration // Retry to time after failed

	// The last request failed
	Failed bool
	Active bool

	InactiveTries int
	ActiveTries   int
	Score         int
}

type Backends []*Backend

type ByScore []*Backend

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].Score < a[j].Score }

func NewBackend(name, address, heartbeat, hbmethod string,
	activeAfter, inactiveAfter, heartbeatTime, retryTime int) *Backend {
	return &Backend{
		Name:      name,
		Address:   address,
		Heartbeat: heartbeat,
		HBMethod:  hbmethod,

		ActiveAfter:   activeAfter,
		InactiveAfter: inactiveAfter,
		HeartbeatTime: time.Duration(heartbeatTime) * time.Millisecond,
		RetryTime:     time.Duration(retryTime) * time.Millisecond,

		Failed:        true,
		Active:        true,
		InactiveTries: 0,
		Score:         0,
	}
}

// Monitoring the backend, can add or remove if heartbeat fail
func (b *Backend) HeartCheck() {
	go func() {
		for {
			var request *http.Request
			var err error

			client := &http.Client{}
			request, err = http.NewRequest(b.HBMethod, b.Heartbeat, nil)
			request.Header.Set("User-Agent", "SSLB-Heartbeat")

			resp, err := client.Do(request)
			if err != nil || resp.StatusCode >= 400 {
				b.RWMutex.Lock()
				// Max tries before consider inactive
				if b.InactiveTries >= b.InactiveAfter {
					log.Printf("Backend inactive [%s]", b.Name)
					b.Active = false
					b.ActiveTries = 0
				} else {
					// Ok that guy it's out of the game
					b.Failed = true
					b.InactiveTries++
					log.Printf("Error to check address [%s] name [%s] tries [%d]", b.Heartbeat, b.Name, b.InactiveTries)
				}
				b.RWMutex.Unlock()
			} else {
				defer resp.Body.Close()

				// Ok, let's keep working boys
				b.RWMutex.Lock()
				if b.ActiveTries >= b.ActiveAfter {
					if b.Failed {
						log.Printf("Backend active [%s]", b.Name)
					}

					b.Failed = false
					b.Active = true
					b.InactiveTries = 0
				} else {
					b.ActiveTries++
				}
				b.RWMutex.Unlock()
			}

			if b.Failed {
				time.Sleep(b.RetryTime)
			} else {
				time.Sleep(b.HeartbeatTime)
			}
		}
	}()
}
