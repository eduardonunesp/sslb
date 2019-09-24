package impl

import (
	"log"

	"github.com/eduardonunesp/sslb/types"
)

type FrontendFactory struct {
	backendFactory        types.BackendFactoryCreator
	backendManagerFactory types.BackendManagerFactoryCreator
}

func NewFrontendFactory() *FrontendFactory {
	return &FrontendFactory{}
}

func (fe FrontendFactory) CreateNewFrontend() types.FrontendRunner {
	if fe.backendFactory == nil {
		log.Fatalln("backendFactory cannot be nil on FrontendFactory")
	}

	if fe.backendManagerFactory == nil {
		log.Fatalln("backendFactory cannot be nil on FrontendFactory")
	}

	f := NewFrontend()
	fe.backendManagerFactory.SetBackendFactory(fe.backendFactory)
	f.SetBackendManagerFactory(fe.backendManagerFactory)
	return f
}

func (fe *FrontendFactory) SetBackendManagerFactory(backendManagerFactory types.BackendManagerFactoryCreator) {
	fe.backendManagerFactory = backendManagerFactory
}

func (fe *FrontendFactory) SetBackendFactory(backendFactory types.BackendFactoryCreator) {
	fe.backendFactory = backendFactory
}
