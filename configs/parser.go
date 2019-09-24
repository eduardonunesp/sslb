package configs

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"runtime"
	"time"

	"github.com/eduardonunesp/sslb/types"
)

const DefaultFilename = "config.json"

// Parse will parse the configuration files and return the types.Config struct
func Parse(filename string) types.Config {
	file := openFile(filename)
	return confParser(file)
}

func confParser(file []byte) types.Config {
	if err := validate(file); err != nil {
		log.Fatal("Can't validate config.json ", err)
	}

	configuration := types.Config{
		GeneralConfig: types.GeneralConfig{
			MaxProcs:         runtime.NumCPU(),
			WorkerPoolSize:   10,
			GracefulShutdown: true,
			Websocket:        true,
			LogLevel:         "info",
			RPCHost:          "127.0.0.1",
			RPCPort:          42555,
		},
		FrontendConfigList: []types.FrontendConfig{
			{
				Timeout: time.Millisecond * 30000,
				BackendConfigList: []types.BackendConfig{
					{
						types.BackendHTTPConfig{
							HBMethod: "HEAD",
						},
						types.BackendControlConfig{
							ActiveAfter:   1,
							InactiveAfter: 3,
							Weight:        1,
						},
						types.BackendHealthCheckConfig{
							HeartbeatTime: time.Millisecond * 30000,
							RetryTime:     time.Millisecond * 5000,
						},
					},
				},
			},
		},
	}

	err := json.Unmarshal(file, &configuration)

	if err != nil {
		log.Fatal("Error to parse json conf", err.Error())
	}

	return configuration
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

	file, err = ioutil.ReadFile("/etc/sslb/" + DefaultFilename)
	if err == nil {
		return file
	}

	file, err = ioutil.ReadFile("~/./sslb/" + DefaultFilename)
	if err == nil {
		return file
	}

	file, err = ioutil.ReadFile("./" + DefaultFilename)
	if err != nil {
		log.Fatal("No config file found, in /etc/sslb or ~/.sslb or in current dir")
	}

	return file
}
