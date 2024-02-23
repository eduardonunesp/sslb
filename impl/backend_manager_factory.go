package impl

import (
	"log"

	"github.com/eduardonunesp/sslb/types"
)

type BackendManagerFactory struct {
	backendFactory types.BackendFactoryCreator
}

func NewBackendManagerFactory() *BackendManagerFactory {
	return &BackendManagerFactory{}
}

func (bmf BackendManagerFactory) CreateNewBackendManager() types.BackendManagerRequester {
	if bmf.backendFactory == nil {
		log.Fatalln("BackendFactory cannot be nil on BackendManagerFactory")
	}

	bf := NewBackendManager()
	bf.SetBackendFactory(bmf.backendFactory)
	return bf
}

func (bmf *BackendManagerFactory) SetBackendFactory(backendFactory types.BackendFactoryCreator) {
	bmf.backendFactory = backendFactory
}
