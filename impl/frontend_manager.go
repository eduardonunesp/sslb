package impl

import (
	"log"

	"github.com/eduardonunesp/sslb/types"
)

type FrontendManager struct {
	frontendConfigList types.FrontendConfigList
	frontendList       types.FrontendRunngerList
	frontendFactory    types.FrontendFactoryCreator
}

func NewFrontendManager() *FrontendManager {
	return &FrontendManager{}
}

func (fe *FrontendManager) SetConfig(config types.FrontendConfigList) {
	fe.frontendConfigList = config
}

func (fe *FrontendManager) SetFrontendFactory(frontendFactory types.FrontendFactoryCreator) {
	fe.frontendFactory = frontendFactory
}

func (fe *FrontendManager) LoadConfig() {
	log.Println("Loading frontend config")

	for _, frontendConfig := range fe.frontendConfigList {
		log.Printf("Loading config for frontend %s\n", frontendConfig.Name)
		newFrontend := fe.frontendFactory.CreateNewFrontend()
		newFrontend.SetConfig(frontendConfig)
		newFrontend.LoadConfig()

		if err := fe.preChecksBeforeAdd(newFrontend); err != nil {
			log.Fatal(err.Error())
		} else {
			fe.frontendList = append(fe.frontendList, newFrontend)
		}
	}
}

func (fe FrontendManager) Run() {
	for _, frontend := range fe.frontendList {
		go frontend.Run()
	}
}

func (fe *FrontendManager) preChecksBeforeAdd(newFrontend types.FrontendRunner) error {
	for _, frontend := range fe.frontendList {
		if frontend.GetConfig().Route == newFrontend.GetConfig().Route {
			return errRouteExists
		}

		if frontend.GetConfig().Port == newFrontend.GetConfig().Port {
			return errPortExists
		}
	}

	return nil
}
