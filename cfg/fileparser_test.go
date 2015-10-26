package cfg

import (
	"runtime"
	"testing"
	"time"
)

func TestFileparserGeneral(t *testing.T) {
	jsonConf := []byte(`
    {
        "general": {
            "maxProcs": 4,
            "gracefulShutdown": true,
            "logLevel": "info",
            "websocket": true,
            "rpchost": "127.0.0.2",
            "rpcport": 42552
        },
        "frontends" : [
            {
                "name" : "Front1",
                "host" : "127.0.0.1",
                "port" : 9000,
                "route" : "/",
                "timeout" : 5000
            }
        ]
    }`)

	conf := ConfParser(jsonConf)
	if conf.GeneralConfig.MaxProcs != runtime.NumCPU() {
		t.Fatal("MaxProcs is wrong", conf.GeneralConfig.MaxProcs)
	}

	if conf.GeneralConfig.WorkerPoolSize != 10 {
		t.Fatal("WorkerPoolSize is wrong", conf.GeneralConfig.WorkerPoolSize)
	}

	if !conf.GeneralConfig.GracefulShutdown {
		t.Fatal("GracefulShutdown is wrong", conf.GeneralConfig.GracefulShutdown)
	}

	if conf.GeneralConfig.LogLevel != "info" {
		t.Fatal("LogLevel is wrong", conf.GeneralConfig.LogLevel)
	}

	if !conf.GeneralConfig.Websocket {
		t.Fatal("Websocket is wrong", conf.GeneralConfig.Websocket)
	}

	if conf.GeneralConfig.RPCHost != "127.0.0.2" {
		t.Fatal("RPCHost is wrong", conf.GeneralConfig.RPCHost)
	}

	if conf.GeneralConfig.RPCPort != 42552 {
		t.Fatal("RPCPort is wrong", conf.GeneralConfig.RPCPort)
	}
}

func TestFileparserFrontend(t *testing.T) {
	jsonConf := []byte(`
    {
        "general": {
            "maxProcs": 4,
            "workerPoolSize": 1000,
            "gracefulShutdown": true,
            "logLevel": "info",
            "websocket": true,
            "rpchost": "127.0.0.2",
            "rpcport": 42552
        },
        "frontends" : [
            {
                "name" : "Front1",
                "host" : "127.0.0.1",
                "port" : 9000,
                "route" : "/"
            }
        ]
    }`)

	conf := ConfParser(jsonConf)
	if conf.FrontendsConfig[0].Name != "Front1" {
		t.Fatal("Name is wrong", conf.FrontendsConfig[0].Name)
	}

	if conf.FrontendsConfig[0].Host != "127.0.0.1" {
		t.Fatal("Host is wrong", conf.FrontendsConfig[0].Host)
	}

	if conf.FrontendsConfig[0].Port != 9000 {
		t.Fatal("Port is wrong", conf.FrontendsConfig[0].Port)
	}

	if conf.FrontendsConfig[0].Route != "/" {
		t.Fatal("Route is wrong", conf.FrontendsConfig[0].Route)
	}

	timeout := time.Millisecond * 30000
	if conf.FrontendsConfig[0].Timeout != timeout {
		t.Fatal("Timeout is wrong", conf.FrontendsConfig[0].Timeout)
	}
}

func TestFileparserBackend(t *testing.T) {
	jsonConf := []byte(`
    {
        "general": {
            "maxProcs": 4,
            "workerPoolSize": 1000,
            "gracefulShutdown": true,
            "logLevel": "info",
            "websocket": true,
            "host": "127.0.0.1",
            "port": 42555
        },

        "frontends" : [
            {
                "name" : "Front1",
                "host" : "127.0.0.1",
                "port" : 9000,
                "route" : "/",
                "timeout" : 5000,

                "backends" : [
                    {
                        "name" : "Back1",
                        "address" : "http://127.0.0.1:9001",
                        "heartbeat" : "http://127.0.0.1:9001",
                        "hbmethod" : "HEAD",
                        "weigth": 1,
                        "inactiveAfter" : 3,
                        "activeAfter" : 1,
                        "heartbeatTime" : 5000,
                        "retryTime" : 1000
                    },
                    {
                        "name" : "Back2",
                        "address" : "http://127.0.0.1:9002",
                        "heartbeat" : "http://127.0.0.1:9002",
                        "hbmethod" : "HEAD",
                        "weigth": 2,
                        "inactiveAfter" : 3,
                        "activeAfter" : 1,
                        "heartbeatTime" : 5000,
                        "retryTime" : 1000
                    }

                ]
            }
        ]
    }`)

	conf := ConfParser(jsonConf)
	if conf.FrontendsConfig[0].BackendsConfig[0].Name != "Back1" {
		t.Fatal("Name is wrong", conf.FrontendsConfig[0].BackendsConfig[0].Name)
	}
}
