package lru

import (
	"fmt"
	"testing"
)

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

func Test_Set(t *testing.T) {
	c := newcache(1)
	c.Put("a", 1)
	v, ok := c.hash["a"]
	if !ok {
		t.Fatal("put not ok")
	}

	if p, _ := v.Value.(*pair); p.val.(int) != 1 {
		t.Fatalf("put 1 not ok. expect=%v, actual=%v", 1, v.Value)
	}

	c.Put("b", 1)
	_, ok = c.hash["b"]
	if !ok {
		t.Fatal("put b not ok")
	}

	_, ok = c.hash["a"]
	if ok {
		t.Fatal("put out size not ok")
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
		t.Fatalf("set get failed. val invalid. expect=%v, actual=%v ", 1, v)
	}

	c.Put("a", 2)
	val, ok = c.Get("a")
	if !ok {
		t.Fatal("set get 2 failed. not ok ")
	}

	v, _ = val.(int)
	if v != 2 {
		t.Fatalf("set get 2 failed. val invalid. expect=%v, actual=%v ", 2, v)
	}
}

func Test_Recent(t *testing.T) {
	c := NewCache(3)
	c.Put("a", 1)
	c.Put("b", 1)
	c.Put("c", 1)
	a, ok := c.Get("a")
	if !ok {
		t.Fatal("get recent failed. get a")
	}

	c.Put("d", 1)

	b, ok := c.Get("b")
	if ok || fmt.Sprintf("%v", b) != fmt.Sprintf("%v", nil) {
		t.Fatalf("get recent failed. expect=%v, atrual=%v", false, ok)
	}

	vc, ok := c.Get("c")
	if !ok || fmt.Sprintf("%v", vc) != "1" {
		t.Fatal("get recent failed. get c")
	}

	obj, ok := c.Get("d")
	if !ok || fmt.Sprintf("%v", obj) != "1" {
		t.Fatal("get recent failed. get c")
	}
	t.Logf("d:%v, a:%v, b:%v, c:%v", obj, a, b, vc)
}
