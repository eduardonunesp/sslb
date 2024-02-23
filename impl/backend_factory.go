package impl

import "github.com/eduardonunesp/sslb/types"

type BackendFactory struct{}

func NewBackendFactory() *BackendFactory {
	return &BackendFactory{}
}

func (bf BackendFactory) CreateNewBackend() types.BackendRequester {
	return NewBackend()
}
