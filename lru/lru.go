package lru

import "container/list"

type Cache interface {
	Get(key string) (interface{}, bool)
	Put(key string, val interface{})
	Clear()
}

type cache struct {
	size int
	hash map[string]*list.Element
	list *list.List
}

func NewCache(size uint) Cache {
	return &cache{
		size: int(size),
		hash: make(map[string]*list.Element, size),
		list: list.New(),
	}
}

func (c *cache) Get(key string) (interface{}, bool) {

	return nil, false
}

func (c *cache) Put(key string, val interface{}) {
	return
}

func (c *cache) Clear() {
}
