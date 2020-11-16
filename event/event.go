package event

import (
	"container/list"
	"context"
)

type EventEntry interface {
	Name() string
	Data() interface{}
	Target() interface{}
}

type EventEmitter interface {
	AddEventListener(name string, handler Handle) *EventHandle
	RemoveEventListener(name string, handle *EventHandle)
	RemoveAll(name string) bool
	Dispatch(ctx context.Context, name string, data interface{})
}

type Handle func(ctx context.Context, event EventEntry)

type EventHandle struct {
	Handle
	*list.Element
	Name string
}

func (eh *EventHandle) Dispose() {
	eh.Element = nil
	eh.Handle = nil
	eh.Name = ""
}

func NewEventHandle(handle Handle) *EventHandle {
	return &EventHandle{
		Handle: handle,
	}
}

// 不支持并发
type eventEmitter struct {
	m map[string]*list.List
}

func (e *eventEmitter) AddEventListener(name string, h Handle) *EventHandle {
	l, ok := e.m[name]
	if !ok {
		l = list.New()
		e.m[name] = l
	}

	event := &EventHandle{
		Name:   name,
		Handle: h,
	}

	ele := l.PushBack(event)
	event.Element = ele

	return event
}

func (e *eventEmitter) RemoveEventListener(name string, eh *EventHandle) {
	if eh == nil || eh.Element == nil {
		return
	}

	l, ok := e.m[name]
	if !ok {
		return
	}

	val := l.Remove(eh.Element)
	t, _ := val.(*EventHandle)
	t.Dispose()

	eh.Dispose()

	return
}

func (e *eventEmitter) RemoveAll(name string) bool {
	l, ok := e.m[name]
	if !ok {
		return false
	}

	delete(e.m, name)

	f := l.Front()
	if f == nil {
		return true
	}

	// dispose every event handle
	for f != nil {
		fv, _ := f.Value.(*EventHandle)
		if fv == nil {
			continue
		}

		fv.Dispose()
		f = f.Next()
	}

	l.Init()

	return true
}

type eventEntry struct {
	data   interface{}
	name   string
	target interface{}
}

func (entry *eventEntry) Data() interface{} {
	return entry.data
}

func (entry *eventEntry) Name() string {
	return entry.name
}

func (entry *eventEntry) Target() interface{} {
	return entry.target
}

func (e *eventEmitter) Dispatch(ctx context.Context, name string, data interface{}) {
	l, ok := e.m[name]
	if !ok {
		return
	}

	f := l.Front()
	if f == nil {
		return
	}

	for f != nil {
		fv, _ := f.Value.(*EventHandle)
		if fv == nil {
			continue
		}

		// todo object pool here
		ee := &eventEntry{name: name, data: data, target: e}

		fv.Handle(ctx, ee)

		f = f.Next()
	}
}
