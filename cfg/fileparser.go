package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/eduardonunesp/sslb/lb"
)

// General config
type General struct {
	MaxProcs int
}

// Backend config
type Backend struct {
	Name          string
	Address       string
	Heartbeat     string
	InactiveAfter int
	HeartbeatTime int
	RetryTime     int
}

// Frontend config
type Frontend struct {
	Name               string
	Host               string
	Port               int
	Route              string
	Timeout            int
	WorkerPoolSize     int
	DispatcherPoolSize int
	Backends           []Backend
}

// Config structured used to build the server
type Config struct {
	General   General
	Frontends []Frontend
}

// Parse JSON FILE
func ConfParser() Config {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		log.Fatal("File error: %v\n", e)
	}

	var jsonConfig Config
	json.Unmarshal(file, &jsonConfig)

	return jsonConfig
}

// Build everything and let the server run
func Setup() *lb.Server {
	config := ConfParser()

	server := lb.NewServer(config.General.MaxProcs)

	for _, frontendConfig := range config.Frontends {
		frontend := lb.NewFrontend(
			frontendConfig.Name, frontendConfig.Host,
			frontendConfig.Port, frontendConfig.Route, frontendConfig.Timeout,
			frontendConfig.WorkerPoolSize, frontendConfig.DispatcherPoolSize)

		for _, backendConfig := range frontendConfig.Backends {
			backend := lb.NewBackend(backendConfig.Name, backendConfig.Address,
				backendConfig.Heartbeat, backendConfig.InactiveAfter, backendConfig.HeartbeatTime,
				backendConfig.RetryTime)
			frontend.AddBackend(backend)
		}

		server.AddFrontend(frontend)
	}

	return server
}
