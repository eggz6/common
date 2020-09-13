package rsync

import (
	"context"
	"time"
)

type Mutex interface {
	Lock(ctx context.Context, key string, d time.Duration) error
	TryLock(ctx context.Context, key string) bool
	Unlock(ctx context.Context, key string)
}

type Lockable interface {
	SetNX(ctx context.Context, key string, val interface{}, d time.Duration) bool
	Del(ctx context.Context, key string)
}
