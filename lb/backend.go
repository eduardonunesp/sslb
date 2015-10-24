package lb

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// BackendConfig it's the configuration loaded
type BackendConfig struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Heartbeat string `json:"heartbeat"`
	HBMethod  string `json:"hbmethod"`

	ActiveAfter   int `json:"activeAfter"`
	InactiveAfter int `json:"inactiveAfter"` // Consider inactive after max inactiveAfter
	Weight        int `json:"weigth"`

	HeartbeatTime time.Duration `json:"heartbeatTime"` // Heartbeat time if health
	RetryTime     time.Duration `json:"retryTime"`     // Retry to time after failed
}

type BackendsConfig []BackendConfig

// BackendControl keep the control data
type BackendControl struct {
	Failed bool // The last request failed
	Active bool

	InactiveTries int
	ActiveTries   int
	Score         int
}

// Backend structure
type Backend struct {
	BackendConfig
	BackendControl
	sync.RWMutex
}

type Backends []*Backend

func NewBackend(backendConfig BackendConfig) *Backend {
	backendConfig.HeartbeatTime = backendConfig.HeartbeatTime * time.Millisecond
	backendConfig.RetryTime = backendConfig.RetryTime * time.Millisecond

	return &Backend{
		BackendConfig: backendConfig,
		BackendControl: BackendControl{
			true, false,
			0, 0, 0,
		},
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
