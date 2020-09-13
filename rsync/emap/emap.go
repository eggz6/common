package emap

import (
	"hash"
	"sync"

	"hash/fnv"

	"github.com/eggz6/common/rsync"
)

type buket struct {
	mu  sync.RWMutex
	buf map[string]interface{}
}

func (b *buket) Buf() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	res := make(map[string]interface{})
	for k, v := range b.buf {
		res[k] = v
	}

	return res
}

func (b *buket) Keys() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	l := len(b.buf)
	res := make([]string, l)
	i := 0

	for k, _ := range b.buf {
		res[i] = k
		i++
	}

	return res
}

func (b *buket) Get(key string) (interface{}, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	val, ok := b.buf[key]

	return val, ok
}

func (b *buket) Put(key string, val interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buf[key] = val
}

func (b *buket) GetAndDo(key string, ac rsync.Action) (interface{}, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	val, ok := b.buf[key]
	if ac != nil {
		v := ac(key, val)
		b.buf[key] = v
	}

	return val, ok
}

type EMap struct {
	bukets []*buket
	seed   uint32
	pool   sync.Pool
	mu     sync.Mutex
}

func NewMap() *EMap {
	const seed = 12
	bs := make([]*buket, seed)

	for i := 0; i < seed; i++ {
		bs[i] = &buket{buf: make(map[string]interface{}, 0)}
	}

	var p sync.Pool
	p.New = func() interface{} {
		return fnv.New32a()
	}

	return &EMap{bukets: bs, seed: seed, pool: p}
}

func (e *EMap) hashFunc(key string) uint32 {
	h := e.pool.Get().(hash.Hash32)
	h.Reset()
	h.Write([]byte(key))
	sum := h.Sum32()
	e.pool.Put(h)

	return sum % e.seed
}

func (e *EMap) Get(key string) (interface{}, bool) {
	idx := e.hashFunc(key)

	return e.bukets[idx].Get(key)
}

func (e *EMap) Put(key string, val interface{}) {
	idx := e.hashFunc(key)

	e.bukets[idx].Put(key, val)
}

func (e *EMap) GetAndDo(key string, ac rsync.Action) (interface{}, bool) {
	idx := e.hashFunc(key)

	return e.bukets[idx].GetAndDo(key, ac)
}

func (e *EMap) KV() map[string]interface{} {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make(map[string]interface{})
	for _, buket := range e.bukets {
		for k, v := range buket.Buf() {
			res[k] = v
		}
	}

	return res
}

func (e *EMap) Keys() [][]string {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([][]string, e.seed)
	for i, buket := range e.bukets {
		res[i] = buket.Keys()
	}

	return res
}
