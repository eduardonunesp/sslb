package types

// Base interfaces

type ConfigSetter interface {
	SetConfig(config Config)
}

type ConfigGetter interface {
	GetConfig() Config
}

type ConfigLoader interface {
	LoadConfig()
}

type ConfigParser interface {
	Parse(filename string) Config
}

// Base structs

type GeneralConfig struct {
	MaxProcs         int    `json:"maxProcs"`
	WorkerPoolSize   int    `json:"workerPoolSize"`
	GracefulShutdown bool   `json:"gracefulShutdown"`
	Websocket        bool   `json:"websocket"`
	LogLevel         string `json:"logLevel"` // Need to define how it works
	RPCHost          string `json:"rpchost"`
	RPCPort          int    `json:"rpcport"`
}

// Composite structs

type Config struct {
	GeneralConfig      `json:"general"`
	FrontendConfigList `json:"frontends"`
}
