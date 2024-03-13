package lru

import (
	"container/list"
	"sync"
)

type Cache interface {
	Get(key string) (interface{}, bool)
	Put(key string, val interface{})
	Clear()
}

type pair struct {
	key string
	val interface{}
}

type cache struct {
	size int
	hash map[string]*list.Element
	list *list.List
	pool sync.Pool

	mux sync.RWMutex
}

func NewCache(size uint) Cache {
	return newcache(size)
}

func newcache(size uint) *cache {
	res := &cache{
		size: int(size),
		hash: make(map[string]*list.Element, size),
		list: list.New(),
	}

	res.pool.New = func() interface{} {
		return &pair{}
	}

	return res
}

func (c *cache) Get(key string) (interface{}, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	v, ok := c.hash[key]
	if !ok {
		return nil, ok
	}

	c.list.MoveToBack(v)

	p, _ := v.Value.(*pair)

	return p.val, ok
}

func (c *cache) allocPair(key string, val interface{}) *pair {
	obj := c.pool.Get()
	kv, _ := obj.(*pair)
	kv.key = key
	kv.val = val

	return kv
}

func (c *cache) Put(key string, val interface{}) {
	c.mux.Lock()
	defer c.mux.Unlock()

	v, ok := c.hash[key]
	if !ok {
		var kv *pair
		if c.list.Len() >= c.size {
			obj := c.list.Remove(c.list.Front())
			kv, _ = obj.(*pair)
			delete(c.hash, kv.key)
			kv.key = ""
			kv.val = nil
		}

		if kv == nil {
			kv = c.allocPair(key, val)
		} else {
			kv.key = key
			kv.val = val
		}

		ele := c.list.PushBack(kv)
		c.hash[key] = ele

		return
	}

	tmp := v.Value
	kv, _ := tmp.(*pair)
	kv.val = val

	return
}

func (c *cache) Clear() {
}
