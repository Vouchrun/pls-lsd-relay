package utils

import "sync"

type KeyedMutex[T comparable] struct {
	mutexes sync.Map
}

func (m *KeyedMutex[T]) Lock(key T) func() {
	value, _ := m.mutexes.LoadOrStore(key, &sync.Mutex{})
	mtx := value.(*sync.Mutex)
	mtx.Lock()

	return func() { mtx.Unlock() }
}
