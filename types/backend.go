package types

import (
	"net/http"
	"time"
)

// Base interfaces

type BackendConfigSetter interface {
	SetConfig(config BackendConfig)
}

type BackengConfigGetter interface {
	GetConfig() BackendConfig
}

type BackendConfigListSetter interface {
	SetConfig(config BackendConfigList)
}

type BackendManagerSetter interface {
	SetBackendManager(backendManager BackendManagerRequester)
}

type BackendManagerGetter interface {
	GetBackendManager() BackendManagerRequester
}

type BackendManagerFactorySetter interface {
	SetBackendManagerFactory(backendManagerFactory BackendManagerFactoryCreator)
}

type BackendManagerFactoryGetter interface {
	GetBackendManagerFactory() BackendManagerFactoryCreator
}

type BackendFactorySetter interface {
	SetBackendFactory(backendFactory BackendFactoryCreator)
}

type BackendFactoryGetter interface {
	GetBackendFactory() BackendFactoryCreator
}

type BackendFactoryCreator interface {
	CreateNewBackend() BackendRequester
}

// Composite interfaces

type BackendManagerFactoryCreator interface {
	BackendFactorySetter
	CreateNewBackendManager() BackendManagerRequester
}

type BackendRequester interface {
	BackendConfigSetter
	BackengConfigGetter
	GetScore() int
	HealthCheck()
	CreateInternalRequest(frontendRequest *http.Request) chan http.Response
}

type BackendManagerRequester interface {
	BackendConfigListSetter
	BackendFactorySetter
	ConfigLoader
	HealthCheck()
	NewRequest(frontendRequest *http.Request) chan http.Response
}

// Base structs

type BackendHTTPConfig struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Heartbeat string `json:"heartbeat"`
	HBMethod  string `json:"hbmethod"`
}

type BackendControlConfig struct {
	ActiveAfter   int `json:"activeAfter"`
	InactiveAfter int `json:"inactiveAfter"` // Consider inactive after max inactiveAfter
	Weight        int `json:"weigth"`
}

type BackendHealthCheckConfig struct {
	HeartbeatTime time.Duration `json:"heartbeatTime"` // Heartbeat time if health
	RetryTime     time.Duration `json:"retryTime"`     // Retry to time after failed
}

// Composite structs

type BackendConfig struct {
	BackendHTTPConfig
	BackendControlConfig
	BackendHealthCheckConfig
}

// List types

type BackendConfigList []BackendConfig
type BackendRequesterList []BackendRequester
