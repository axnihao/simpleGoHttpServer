package factory

import (
	"bookstore/store"
	"fmt"
	"sync"
)

var (
	providerMu sync.RWMutex
	providers  = make(map[string]store.Store)
)

func Register(name string, p store.Store) {
	providerMu.Lock()
	defer providerMu.Unlock()
	if p == nil {
		panic("store: Register provider is nil")
	}

	if _, dup := providers[name]; dup {
		panic("store: Register called twice for provider " + name)
	}
	providers[name] = p
}

func New(providerName string) (store.Store, error) {
	providerMu.RLock()
	p, ok := providers[providerName]
	providerMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("store: unknown provider %s", providerName)
	}
	return p, nil

}
