package keypool

import (
	"data-lib/uniqueboundedmap"
	"sync"
)

type KeyPool[K comparable] struct {
	m       *uniqueboundedmap.UniqueBoundedMap[K, chan struct{}]
	chMap   map[K]chan struct{}
	mapLock sync.RWMutex
}

func NewKeyPool[K comparable](maxSize int) *KeyPool[K] {
	return &KeyPool[K]{
		m: uniqueboundedmap.NewUniqueBoundedMap[K, chan struct{}](maxSize),
	}
}

func (p *KeyPool[K]) Acquire(key K) {
	for {
		if ch, ok := p.m.Get(key); ok {
			<-ch
			continue
		}
		ch := make(chan struct{})
		if res := p.m.Put(key, ch); res != uniqueboundedmap.OK {
			continue
		}
		p.mapLock.Lock()
		p.chMap[key] = ch
		p.mapLock.Unlock()
		break
	}
}

func (p *KeyPool[K]) Release(key K) {
	p.mapLock.Lock()
	ch := p.chMap[key]
	delete(p.chMap, key)
	p.mapLock.Unlock()
	close(ch)
	p.m.Delete(key)
}
