package types

// Base interfaces

type ServerStopper interface {
	Stop()
}

// Composite interfaces

type ServerRunner interface {
	ServerStopper
	ConfigSetter
	FrontendManagerFactorySetter
	Run()
}
