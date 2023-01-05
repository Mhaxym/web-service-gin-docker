package goCache

type GoCacheManager struct {
	elements         map[string]IManager
	blockingChannels map[string]chan struct{}
}

type IGoCache interface {
	Get(managerName string) IManager
	Set(value IManager) IManager
}

type IManager interface {
	Load()
	GetManagerName() string 
}

type Manager struct {
	CacheName string
}
