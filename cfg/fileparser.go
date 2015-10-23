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
	HBMethod      string
	ActiveAfter   int
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

const DEFAULT_FILENAME = "config.json"

// Config structured used to build the server
type Config struct {
	General   General
	Frontends []Frontend
}

func openFile(filename string) []byte {
	var file []byte
	var err error

	if filename != "" {
		file, err = ioutil.ReadFile(filename)
		if err == nil {
			return file
		} else {
			log.Fatal(err)
		}
	}

	file, err = ioutil.ReadFile("/etc/sslb/" + DEFAULT_FILENAME)
	if err == nil {
		return file
	}

	file, err = ioutil.ReadFile("~/./sslb/" + DEFAULT_FILENAME)
	if err == nil {
		return file
	}

	file, err = ioutil.ReadFile("./" + DEFAULT_FILENAME)
	if err != nil {
		log.Fatal("No config file found, in /etc/sslb or ~/.sslb or in current dir")
	}

	return file
}

// ConfParser to Parse JSON FILE
func ConfParser(filename string) Config {
	file := openFile(filename)

	var jsonConfig Config
	err := json.Unmarshal(file, &jsonConfig)

	if err != nil {
		log.Fatal("Error to parse json conf", err.Error())
	}

	return jsonConfig
}

func CreateConfig(filename string) {
	configExample := []byte(`{
    "general": {
        "maxProcs": 4,
        "workerPoolSize": 1000,
        "dispatcherPoolSize": 1000
    },

    "frontends" : [{
        "name" : "Frontend App",
        "host" : "0.0.0.0",
        "port" : 80,
        "route" : "/",
        "timeout" : 5000,

        "backends" : [
            {
                "name" : "Backend 1",
                "address" : "http://127.0.0.1:9001",
                "heartbeat" : "http://127.0.0.1:9001/heartbeat",
                "inactiveAfter" : 3,
                "heartbeatTime" : 15000,
                "retryTime" : 1000
            },{
                "name" : "Backend 2",
                "address" : "http://127.0.0.1:9002",
                "heartbeat" : "http://127.0.0.1:9002/heartbeat",
                "hbmethod" : "HEAD",
                "inactiveAfter" : 3,
                "activeAfter" : 1,
                "heartbeatTime" : 15000,
                "retryTime" : 1000
            },{
                "name" : "Backend 3",
                "address" : "http://127.0.0.1:9003",
                "heartbeat" : "http://127.0.0.1:9003/heartbeat",
                "hbmethod" : "HEAD",
                "activeAfter" : 1,
                "inactiveAfter" : 1,
                "heartbeatTime" : 5000,
                "retryTime" : 1000
            }
        ]
    }]
}`)

	err := ioutil.WriteFile(filename, configExample, 0644)
	if err != nil {
		log.Fatal("Can't create file config.json example", err)
	}
}

// Setup will build everything and let the server run
func Setup(filename string) *lb.Server {
	config := ConfParser(filename)

	cpus := runtime.NumCPU()
	log.Printf("%d CPUS available, using only %d", cpus, config.General.MaxProcs)

	runtime.GOMAXPROCS(config.General.MaxProcs)

	server := lb.NewServer(config.General.WorkerPoolSize, config.General.DispatcherPoolSize)

	for _, frontendConfig := range config.Frontends {
		frontend := endpoint.NewFrontend(
			frontendConfig.Name, frontendConfig.Host,
			frontendConfig.Port, frontendConfig.Route, frontendConfig.Timeout)

		for _, backendConfig := range frontendConfig.Backends {
			backend := endpoint.NewBackend(backendConfig.Name, backendConfig.Address, backendConfig.Heartbeat,
				backendConfig.HBMethod, backendConfig.ActiveAfter, backendConfig.InactiveAfter, backendConfig.HeartbeatTime,
				backendConfig.RetryTime)
			frontend.AddBackend(backend)
		}

		server.AddFrontend(frontend)
	}

	return server
}
