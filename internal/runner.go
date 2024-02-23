package internal

import (
	"log"

	"github.com/eduardonunesp/sslb/types"
)

func RunServer(
	server types.ServerRunner,
	frontendManagerFactory types.FrontendManagerFactoryCreator,
	backendManagerFactory types.BackendManagerFactoryCreator,
	frontendFactory types.FrontendFactoryCreator,
	backendFactory types.BackendFactoryCreator,
) {
	log.Println("Initializing SSLB")

	frontendManagerFactory.SetBackendFactory(backendFactory)
	frontendManagerFactory.SetFrontendFactory(frontendFactory)
	frontendManagerFactory.SetBackendManagerFactory(backendManagerFactory)

	server.SetFrontendManagerFactory(frontendManagerFactory)

	server.Run()
}
