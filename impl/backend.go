package impl

import (
	"log"
	"net/http"
	"time"

	"github.com/eduardonunesp/sslb/types"
)

type Backend struct {
	types.BackendConfig

	Failed bool // The last request failed
	Active bool

	InactiveTries int
	ActiveTries   int
	Score         int
}

type BackendList []Backend

func NewBackend() *Backend {
	return &Backend{}
}

func (b Backend) GetScore() int {
	return b.Score
}

func (b *Backend) SetConfig(config types.BackendConfig) {
	b.BackendConfig = config
}

func (b Backend) GetConfig() types.BackendConfig {
	return b.BackendConfig
}

func (b *Backend) CreateInternalRequest(frontendRequest *http.Request) chan http.Response {
	iRequestChan := make(chan http.Response)

	go func() {
		requestAddress := b.BackendConfig.Address + frontendRequest.URL.String()
		log.Printf("Request on backend %s\n", requestAddress)

		client := &http.Client{}
		httpRequest, _ := http.NewRequest(frontendRequest.Method, requestAddress, frontendRequest.Body)

		for k, vv := range frontendRequest.Header {
			for _, v := range vv {
				httpRequest.Header.Set(k, v)
			}
		}

		response, err := client.Do(httpRequest)

		if err != nil {
			log.Printf("Error on request %v\n", err)
			iRequestChan <- http.Response{StatusCode: 502, Body: nil, Header: frontendRequest.Header}
			return
		}

		iRequestChan <- *response
	}()

	return iRequestChan
}

func (b *Backend) HealthCheck() {
	go func() {
		for {
			var request *http.Request
			var err error

			client := &http.Client{}
			request, err = http.NewRequest(b.HBMethod, b.Heartbeat, nil)
			request.Header.Set("User-Agent", "SSLB-Heartbeat")

			resp, err := client.Do(request)
			if err != nil || resp.StatusCode >= 400 {

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
			} else {
				defer resp.Body.Close()
				// Ok, let's keep working boys

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
			}

			if b.Failed {
				time.Sleep(b.RetryTime * time.Millisecond)
			} else {
				time.Sleep(b.HeartbeatTime * time.Millisecond)
			}
		}
	}()
}
