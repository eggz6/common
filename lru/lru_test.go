package lru

import "testing"

func Test_New(t *testing.T) {
	c := NewCache(10)
	if c == nil {
		t.Fatal("new cache failed.")
	}
}

func Test_Get(t *testing.T) {
	c := NewCache(10)
	_, ok := c.Get("a")
	if ok {
		t.Fatal("get failed. ")
	}
}

func Test_SetGet(t *testing.T) {
	c := NewCache(10)
	c.Put("a", 1)
	val, ok := c.Get("a")
	if !ok {
		t.Fatal("set get failed. not ok ")
	}

	v, _ := val.(int)
	if v != 1 {
		t.Fatal("set get failed. val invalid. expect=%v, actual=%v ", 1, v)
	}
}
