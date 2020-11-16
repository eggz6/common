package event

import (
	"container/list"
	"context"
	"fmt"
	"testing"
)

func Test_AddEvent(t *testing.T) {
	emitter := &eventEmitter{
		m: make(map[string]*list.List, 0),
	}

	emitter.AddEventListener("a", func(ctx context.Context, e EventEntry) {
		fmt.Println("a")
	})

	emitter.AddEventListener("a", func(ctx context.Context, e EventEntry) {
		fmt.Println("a1")
	})

	emitter.AddEventListener("a", func(ctx context.Context, e EventEntry) {
		fmt.Println("a2")
	})

	l, ok := emitter.m["a"]
	if !ok {
		t.Fatalf("no list")
	}

	f := l.Front()
	if f == nil {
		t.Fatalf("no ele")
	}

	i := 0
	for f != nil {
		h, ok := f.Value.(*EventHandle)
		if !ok {
			t.Fatalf("no handle %T", f.Value)
		}

		t.Logf("handle=%v, ", h.Name)
		if h.Name != "a" {
			t.Fatalf("failed name handle=%v, ", h.Name)
		}

		f = h.Next()
		i++
	}
}

func Test_RemoveEvent(t *testing.T) {
	emitter := &eventEmitter{
		m: make(map[string]*list.List, 0),
	}

	emitter.RemoveEventListener("b", &EventHandle{})
	emitter.RemoveEventListener("b", &EventHandle{Element: &list.Element{}})

	a := emitter.AddEventListener("a", func(ctx context.Context, e EventEntry) {
		fmt.Println("a")
	})

	a1 := emitter.AddEventListener("a", func(ctx context.Context, e EventEntry) {
		fmt.Println("a1")
	})

	emitter.RemoveEventListener(a.Name, a)

	l, _ := emitter.m["a"]
	f, ok := l.Front().Value.(*EventHandle)

	if !ok {
		t.Fatalf("remove not a event ")
	}

	if f != a1 {
		t.Fatalf("remove node a event ")
	}

	emitter.RemoveEventListener(a1.Name, a1)
	if l.Len() != 0 {
		t.Fatalf("no elements")
	}

	if a1.Element != nil || a1.Handle != nil || a1.Name != "" {
		t.Fatalf("dispose failed. a1=%v", a)
	}
}

func Test_RemoveAll(t *testing.T) {
	emitter := &eventEmitter{
		m: make(map[string]*list.List, 0),
	}

	emitter.AddEventListener("a", func(ctx context.Context, e EventEntry) { fmt.Println("a") })
	emitter.AddEventListener("a", func(ctx context.Context, e EventEntry) { fmt.Println("a1") })

	ok := emitter.RemoveAll("a")
	if !ok {
		t.Fatalf("remove all failed")
	}

	_, has := emitter.m["a"]

	if has {
		t.Fatalf("no remove")
	}
}

func Test_Dispatch(t *testing.T) {
	emitter := &eventEmitter{
		m: make(map[string]*list.List, 0),
	}

	count := 0
	emitter.AddEventListener("a", func(ctx context.Context, e EventEntry) {
		if e.Name() != "a" {
			t.Fatalf("disptach name failed. want=%v, actural=%v", "a", e.Name())
		}

		if e.Data() == nil {
			t.Fatal("disptach data failed")
		}

		val, _ := e.Data().(int)
		if val != 1 {
			t.Fatalf("disptach val failed. want=%v, actural=%v", 1, e.Data())
		}

		if count != 0 {
			t.Fatalf("disptach count failed. want=%v, actural=%v", 0, count)
		}

		count = count + val
	})

	emitter.AddEventListener("a", func(ctx context.Context, e EventEntry) {
		if e.Name() != "a" {
			t.Fatalf("disptach name failed. want=%v, actural=%v", "a", e.Name())
		}

		if e.Data() == nil {
			t.Fatal("disptach data failed")
		}

		val, _ := e.Data().(int)
		if val != 1 {
			t.Fatalf("disptach val failed. want=%v, actural=%v", 1, e.Data())
		}

		if count != 1 {
			t.Fatalf("disptach count failed. want=%v, actural=%v", 0, count)
		}

		count = count + val
	})

	emitter.Dispatch(context.Background(), "a", 1)

	if count != 2 {
		t.Fatal("no handle run")
	}
}

func Test_Dispatch2(t *testing.T) {
	emitter := &eventEmitter{
		m: make(map[string]*list.List, 0),
	}

	count := 0
	emitter.AddEventListener("a", func(ctx context.Context, e EventEntry) {

		count = count + 1
	})

	emitter.Dispatch(context.Background(), "b", 1)

	if count == 1 {
		t.Fatal("event name error")
	}
}
