package uniqueboundedmap

import (
	"fmt"
	"sync"
)

type PutResult int

const (
	OK PutResult = iota
	KeyExists
)

type UniqueBoundedMap[K comparable, V any] struct {
	m     map[K]V
	mutex sync.RWMutex
	cond  *sync.Cond
	size  int
}

func NewUniqueBoundedMap[K comparable, V any](maxSize int) *UniqueBoundedMap[K, V] {
	l := &UniqueBoundedMap[K, V]{
		m:    make(map[K]V),
		size: maxSize,
	}
	l.cond = sync.NewCond(&l.mutex)
	return l
}

func (l *UniqueBoundedMap[K, V]) Put(key K, value V) PutResult {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for len(l.m) >= l.size {
		l.cond.Wait()
	}

	if _, ok := l.m[key]; ok {
		return KeyExists
	}

	// Add the item to the map
	l.m[key] = value
	fmt.Printf("Added: %s -> %d\n", key, value)
	return OK
}

func (l *UniqueBoundedMap[K, V]) Get(key K) (V, bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	val, ok := l.m[key]
	return val, ok
}

func (l *UniqueBoundedMap[K, V]) Delete(key K) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if _, ok := l.m[key]; ok {
		delete(l.m, key)
		fmt.Printf("Deleted: %s\n", key)
		l.cond.Signal() // Signal one goroutine that there is space
	}
}
