package impl

import (
	"log"
	"net/http"
	"sort"
	"sync"

	"github.com/eduardonunesp/sslb/types"
)

type BackendManager struct {
	backendConfigs types.BackendConfigList
	backendList    types.BackendRequesterList
	backendFactory types.BackendFactoryCreator
	sync.RWMutex
}

func NewBackendManager() *BackendManager {
	return &BackendManager{}
}

func (be *BackendManager) SetConfig(configs types.BackendConfigList) {
	be.backendConfigs = configs
}

func (be *BackendManager) SetBackendFactory(backendFactroy types.BackendFactoryCreator) {
	be.backendFactory = backendFactroy
}

func (be *BackendManager) HealthCheck() {
	for _, backend := range be.backendList {
		log.Printf("Check health on backend%s\n", backend.GetConfig().Name)
		backend.HealthCheck()
	}
}

func (be *BackendManager) getLessScored() types.BackendRequester {
	be.RWMutex.RLock()
	defer be.RWMutex.RUnlock()

	bList := be.backendList[:]

	sort.SliceStable(bList, func(i, j int) bool {
		return bList[i].GetScore() < bList[j].GetScore()
	})

	return bList[0]
}

func (be *BackendManager) NewRequest(frontendRequest *http.Request) chan http.Response {
	return be.getLessScored().CreateInternalRequest(frontendRequest)
}

func (be *BackendManager) LoadConfig() {
	log.Println("Loading backend config")

	if be.backendFactory == nil {
		log.Fatalln("backendFactory is nil on BackendManager")
	}

	for _, backendConfig := range be.backendConfigs {
		log.Printf("Load backend config %s\n", backendConfig.Name)
		newBackend := be.backendFactory.CreateNewBackend()
		newBackend.SetConfig(backendConfig)
		be.backendList = append(be.backendList, newBackend)
	}
}
