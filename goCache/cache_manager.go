package goCache

import "sync"

var goCacheSingleton sync.Once
var goChannelMutex sync.Mutex
var goCacheManager *GoCacheManager

// GetGoCacheManager returns a singleton instance of GoCacheManager.
// This is used to access the GoCacheManager's functions.
func GetGoCacheManager() *GoCacheManager {
	goCacheSingleton.Do(func() {
		goCacheManager = &GoCacheManager{elements: make(map[string]IManager), blockingChannels: make(map[string]chan struct{})}
	})
	return goCacheManager
}

// GoCacheManager.Get returns the IManager, identified by managerName, from the gcm.elements map
// If the IManager does not exist, it returns nil
func (gcm *GoCacheManager) Get(managerName string) IManager {
	if IManager, exists := gcm.elements[managerName]; exists {
		return IManager
	}
	return nil
}

func (gcm *GoCacheManager) Set(value IManager) IManager {
	managerName := value.GetManagerName()
	if gcm.elements[managerName] == nil {
		gcm.lock(managerName)
		if gcm.elements[managerName] == nil {
			value.Load()
			gcm.elements[managerName] = value
		}
		gcm.unlock(managerName)
	}
	return gcm.elements[managerName]
}

func (gcm *GoCacheManager) lock(managerName string) {
	// Check that the channel exists. If not, create it.
	if gcm.blockingChannels[managerName] == nil {
		goChannelMutex.Lock()
		defer goChannelMutex.Unlock()
		if gcm.blockingChannels[managerName] == nil {
			gcm.blockingChannels[managerName] = make(chan struct{}, 1)
			gcm.blockingChannels[managerName] <- struct{}{}
		}

	}
	// Wait until the channel is available
	<-gcm.blockingChannels[managerName]
}

// unlock releases the lock on the manager with the given name.
// This function is used to send a message to the lock manager goroutine
// to tell it to unlock the manager with the given name.
func (gcm *GoCacheManager) unlock(managerName string) {
	gcm.blockingChannels[managerName] <- struct{}{}
}
