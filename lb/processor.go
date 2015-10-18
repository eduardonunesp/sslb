package lb

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

var (
	errTimeout        = errors.New("Timeout")
	errInvalidRequest = errors.New("Invalid request")
)

func Run() {
	runtime.GOMAXPROCS(1)

	host := "localhost"
	port := 9000
	address := fmt.Sprintf("%s:%d", host, port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		chanResponse := WorkerRun(r)
		defer close(chanResponse)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		select {
		case result := <-chanResponse:
			if result.Status > 400 {
				http.Error(w, string(result.Body), result.Status)
			} else {
				w.WriteHeader(result.Status)
				w.Write(result.Body)
			}

		case <-ticker.C:
			http.Error(w, errTimeout.Error(), http.StatusRequestTimeout)
		}
	})

	http.ListenAndServe(address, nil)
}
