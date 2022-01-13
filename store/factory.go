package store

import (
	"fmt"
	"sync"
)

var (
	providersLock = sync.RWMutex{} // 读写锁
	providers =make(map[string]Store)
)

// Register 注册存储模块
func Register(name string, provider Store) {
	if provider == nil {
		panic("store: Register provider is nil")
	}
	providersLock.RLock()
	if _, ok := providers[name]; ok {
		providersLock.RUnlock()
		panic("store: Register called twice for provider " + name)
	}
	providersLock.RUnlock()

	providersLock.Lock()
	defer providersLock.Unlock()
	providers[name] = provider
}

// New 获取存储模块
func New(name string) (store Store, err error) {
	providersLock.RLock()
	store,exist := providers[name]
	providersLock.RUnlock()
	if !exist {
		err = fmt.Errorf("store: unknown provider %s",  name)
		return
	}
	return
}