package emap

import (
	"fmt"
	"testing"
)

func Test_Get(t *testing.T) {
	val := 123
	m := NewMap()
	m.Put("key", val)

	v, ok := m.Get("key")
	if !ok {
		t.Fatal(" has no key")
	}

	if v != val {
		t.Fatal("val no equal")
	}
}

func Test_KV(t *testing.T) {
	m := NewMap()
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key_%v", i)
		m.Put(key, i)
	}

	t.Logf("kv=%v", m.KV())
}

func Test_Keys(t *testing.T) {
	m := NewMap()
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key_%v", i)
		m.Put(key, i)
	}

	keys := m.Keys()
	for _, arr := range keys {
		t.Logf("len=%v, keys=%v", len(arr), arr)
	}
}

func Benchmark_Get(b *testing.B) {
	m := NewMap()
	keys := []string{}
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key_%v", i)
		keys = append(keys, key)
		m.Put(key, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(keys[i%100])
	}
}

func Benchmark_Put(b *testing.B) {
	m := NewMap()
	keys := []string{}
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key_%v", i)
		keys = append(keys, key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := keys[i%100]
		m.Put(key, i)
	}
}
