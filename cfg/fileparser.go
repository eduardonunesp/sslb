package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"runtime"
	"time"

	"github.com/eduardonunesp/sslb/lb"
)

const DEFAULT_FILENAME = "config.json"

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
func ConfParser(file []byte) lb.Configuration {
	if err := Validate(file); err != nil {
		log.Fatal("Can't validate config.json ", err)
	}

	jsonConfig := lb.Configuration{
		GeneralConfig: lb.GeneralConfig{
			MaxProcs:         runtime.NumCPU(),
			WorkerPoolSize:   10,
			GracefulShutdown: true,
			Websocket:        true,
			LogLevel:         "info",
			RPCHost:          "127.0.0.1",
			RPCPort:          42555,
		},
		FrontendsConfig: []lb.FrontendConfig{
			{
				Timeout: time.Millisecond * 30000,
				BackendsConfig: []lb.BackendConfig{
					{
						HBMethod:      "HEAD",
						ActiveAfter:   1,
						InactiveAfter: 3,
						Weight:        1,
						HeartbeatTime: time.Millisecond * 30000,
						RetryTime:     time.Millisecond * 5000,
					},
				},
			},
		},
	}

	err := json.Unmarshal(file, &jsonConfig)

	if err != nil {
		log.Fatal("Error to parse json conf", err.Error())
	}

	return jsonConfig
}

// Setup will build everything and let the server run
func Setup(filename string) lb.Configuration {
	file := openFile(filename)
	return ConfParser(file)
}
