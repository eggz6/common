package common

import "context"

var (
	empty Nil
	key   = meta("metadata")
)

type Nil struct {
}

type meta string

type Metadata map[string]interface{}

func (m Metadata) Get(key string) interface{} {
	val, ok := m[key]
	if !ok {
		return &empty
	}

	return val
}

func (m Metadata) Set(key string, val interface{}) {
	m[key] = val
}

func Meta(ctx context.Context) (Metadata, bool) {
	val := ctx.Value(key)
	if val == nil {
		return nil, false
	}

	res, ok := val.(Metadata)

	return res, ok
}

func Ctx(ctx context.Context) (context.Context, Metadata) {
	m := Metadata(make(map[string]interface{}))
	res := context.WithValue(ctx, key, m)

	return res, m
}
