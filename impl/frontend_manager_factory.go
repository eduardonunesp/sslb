package impl

import (
	"log"

	"github.com/eduardonunesp/sslb/types"
)

type FrontendManagerFactory struct {
	frontendFactory       types.FrontendFactoryCreator
	backendFactory        types.BackendFactoryCreator
	backendManagerFactory types.BackendManagerFactoryCreator
}

func NewFrontendManagerFactory() *FrontendManagerFactory {
	return &FrontendManagerFactory{}
}

func (fmf *FrontendManagerFactory) CreateNewFrontendManager() types.FrontendManagerRunner {
	if fmf.frontendFactory == nil {
		log.Fatalln("frontendFactory cannot be nil on FrontendManagerFactory")
	}

	if fmf.backendFactory == nil {
		log.Fatalln("backendFactory cannot be nil on FrontendManagerFactory")
	}

	if fmf.backendManagerFactory == nil {
		log.Fatalln("backendManagerFactory cannot be nil on FrontendManagerFactory")
	}

	fm := NewFrontendManager()
	fmf.backendManagerFactory.SetBackendFactory(fmf.backendFactory)
	fmf.frontendFactory.SetBackendManagerFactory(fmf.backendManagerFactory)
	fmf.frontendFactory.SetBackendFactory(fmf.backendFactory)
	fm.SetFrontendFactory(fmf.frontendFactory)
	return fm
}

func (fmf *FrontendManagerFactory) SetFrontendFactory(frontendFactory types.FrontendFactoryCreator) {
	fmf.frontendFactory = frontendFactory
}

func (fmf *FrontendManagerFactory) SetBackendFactory(backendFactory types.BackendFactoryCreator) {
	fmf.backendFactory = backendFactory
}

func (fmf *FrontendManagerFactory) SetBackendManagerFactory(backendManagerFactory types.BackendManagerFactoryCreator) {
	fmf.backendManagerFactory = backendManagerFactory
}
