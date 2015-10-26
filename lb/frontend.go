package lb

import (
	"sync"
	"time"
)

// FrontendConfig it's the configuration loaded
type FrontendConfig struct {
	Name           string        `json:"name"`
	Host           string        `json:"host"`
	Port           int           `json:"port"`
	Route          string        `json:"route"`
	Timeout        time.Duration `json:"timeout"`
	BackendsConfig `json:"backends"`
}

type FrontendsConfig []FrontendConfig

// Frontend structure
type Frontend struct {
	FrontendConfig
	Backends
	sync.RWMutex
}

type Frontends []*Frontend

func NewFrontend(frontendConfig FrontendConfig) *Frontend {
	frontendConfig.Timeout = frontendConfig.Timeout * time.Millisecond
	return &Frontend{
		FrontendConfig: frontendConfig,
	}
}
