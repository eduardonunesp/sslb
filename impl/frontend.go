package impl

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/eduardonunesp/sslb/types"
)

type Frontend struct {
	types.FrontendConfig
	backendManagerFactory types.BackendManagerFactoryCreator
	backendManager        types.BackendManagerRequester
}

func NewFrontend() *Frontend {
	return &Frontend{}
}

func (f *Frontend) SetConfig(configuration types.FrontendConfig) {
	f.FrontendConfig = configuration
}

func (f Frontend) GetConfig() types.FrontendConfig {
	return f.FrontendConfig
}

func (f *Frontend) SetBackendManagerFactory(backendManagerFactory types.BackendManagerFactoryCreator) {
	f.backendManagerFactory = backendManagerFactory
}

func (f *Frontend) LoadConfig() {
	f.backendManager = f.backendManagerFactory.CreateNewBackendManager()
	f.backendManager.SetConfig(f.FrontendConfig.BackendConfigList)
	f.backendManager.LoadConfig()
}

func (f *Frontend) Run() {
	// Make sure that frontend has any backend available to request
	if len(f.FrontendConfig.BackendConfigList) == 0 {
		log.Fatal(errNoBackend.Error())
	}

	host := f.FrontendConfig.Host
	port := f.FrontendConfig.Port
	address := fmt.Sprintf("%s:%d", host, port)

	log.Printf("Running from end for %s on address %s\n", f.FrontendConfig.Name, address)

	log.Printf("Start health chech for %s\n", f.FrontendConfig.Name)
	f.backendManager.HealthCheck()

	// Prepare the mux
	httpHandle := http.NewServeMux()

	httpHandle.HandleFunc(f.FrontendConfig.Route, func(w http.ResponseWriter, r *http.Request) {
		// On a serious problem
		defer func() {
			if err := recover(); err != nil {
				log.Printf("%v\n", err)
				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)
			}
		}()

		// Creates a new request for backend returns a channel to consume
		chanResponse := f.backendManager.NewRequest(r)
		r.Close = true

		// Timeout ticker
		ticker := time.NewTicker(f.FrontendConfig.Timeout * time.Millisecond)
		defer ticker.Stop()

		select {
		case result := <-chanResponse:
			// We have a response, it's valid ?
			for k, vv := range result.Header {
				for _, v := range vv {
					w.Header().Set(k, v)
				}
			}

			w.WriteHeader(result.StatusCode)

			if result.Body != nil {
				body, _ := ioutil.ReadAll(result.Body)

				if r.Method != "HEAD" {
					w.Write(body)
				}
			}

		case <-ticker.C:
			// Timeout
			http.Error(w, errTimeout.Error(), http.StatusRequestTimeout)
		}
	})

	// Config and start server
	server := &http.Server{
		Addr:    address,
		Handler: httpHandle,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
