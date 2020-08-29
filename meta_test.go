package common

import (
	"context"
	"testing"
)

func TestMeta(t *testing.T) {
	ctx, meta := Ctx(context.Background())

	meta.Set("123", "val")

	if ctx == nil {
		t.Fatalf("ctx failed")
	}

	m, _ := Meta(ctx)
	if meta.Get("123") != m.Get("123") {
		t.Fatalf("meta failed")
	}
}
