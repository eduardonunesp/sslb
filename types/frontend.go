package types

import (
	"time"
)

// Base interfaces

type FrontendConfigSetter interface {
	SetConfig(frontendConfig FrontendConfig)
}

type FrontendConfigGetter interface {
	GetConfig() FrontendConfig
}

type FrontendFactorySetter interface {
	SetFrontendFactory(frontendFactory FrontendFactoryCreator)
}

type FrontendConfigListSetter interface {
	SetConfig(config FrontendConfigList)
}

type FrontendFactoryGetter interface {
	GetFrontEndFactory() FrontendFactoryCreator
}

type FrontendManagerSetter interface {
	SetFrontendManager(frontendManager FrontendManagerRunner)
}

type FrontendManagerFactorySetter interface {
	SetFrontendManagerFactory(frontendManagerFactory FrontendManagerFactoryCreator)
}

type FrontendFactoryCreator interface {
	BackendManagerFactorySetter
	BackendFactorySetter
	CreateNewFrontend() FrontendRunner
}

// Composite interfaces

type FrontendManagerFactoryCreator interface {
	FrontendFactorySetter
	BackendFactorySetter
	BackendManagerFactorySetter
	CreateNewFrontendManager() FrontendManagerRunner
}

type FrontendRunner interface {
	BackendManagerFactorySetter
	FrontendConfigGetter
	FrontendConfigSetter
	ConfigLoader
	Run()
}

type FrontendManagerRunner interface {
	FrontendConfigListSetter
	ConfigLoader
	Run()
}

// Base structs

type FrontendConfig struct {
	Name              string        `json:"name"`
	Host              string        `json:"host"`
	Port              int           `json:"port"`
	Route             string        `json:"route"`
	Timeout           time.Duration `json:"timeout"`
	BackendConfigList `json:"backends"`
}

// List types

type FrontendConfigList []FrontendConfig
type FrontendRunngerList []FrontendRunner
