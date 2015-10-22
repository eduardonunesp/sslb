package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"runtime"

	"github.com/eduardonunesp/sslb/lb"
	"github.com/eduardonunesp/sslb/lb/endpoint"
)

// General config
type General struct {
	MaxProcs           int
	WorkerPoolSize     int
	DispatcherPoolSize int
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
	Name     string
	Host     string
	Port     int
	Route    string
	Timeout  int
	Backends []Backend
}

// Config structured used to build the server
type Config struct {
	General   General
	Frontends []Frontend
}

// ConfParser to Parse JSON FILE
func ConfParser() Config {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		log.Fatal("File error", e)
	}

	var jsonConfig Config
	err := json.Unmarshal(file, &jsonConfig)

	if err != nil {
		log.Fatal("Error to parse json conf", err.Error())
	}

	return jsonConfig
}

// Setup will build everything and let the server run
func Setup() *lb.Server {
	config := ConfParser()

	cpus := runtime.NumCPU()
	log.Printf("%d CPUS available, using only %d", cpus, config.General.MaxProcs)

	runtime.GOMAXPROCS(config.General.MaxProcs)

	server := lb.NewServer(config.General.WorkerPoolSize, config.General.DispatcherPoolSize)

	for _, frontendConfig := range config.Frontends {
		frontend := endpoint.NewFrontend(
			frontendConfig.Name, frontendConfig.Host,
			frontendConfig.Port, frontendConfig.Route, frontendConfig.Timeout)

		for _, backendConfig := range frontendConfig.Backends {
			backend := endpoint.NewBackend(backendConfig.Name, backendConfig.Address,
				backendConfig.Heartbeat, backendConfig.InactiveAfter, backendConfig.HeartbeatTime,
				backendConfig.RetryTime)
			frontend.AddBackend(backend)
		}

		server.AddFrontend(frontend)
	}

	return server
}
