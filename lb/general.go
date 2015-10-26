package lb

type GeneralConfig struct {
	MaxProcs         int    `json:"maxProcs"`
	WorkerPoolSize   int    `json:"workerPoolSize"`
	GracefulShutdown bool   `json:"gracefulShutdown"`
	Websocket        bool   `json:"websocket"`
	LogLevel         string `json:"logLevel"` // Need to define how it works
	RPCHost          string `json:"rpchost"`
	RPCPort          int    `json:"rpcport"`
}

type Configuration struct {
	GeneralConfig   `json:"general"`
	FrontendsConfig `json:"frontends"`
}
