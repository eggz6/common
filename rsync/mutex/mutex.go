package mutex

import (
	"context"
	"time"

	"github.com/eggz6/common/rsync"
)

type MutexErr string

func (e MutexErr) Error() string {
	return string(e)
}

const (
	LockExpired = MutexErr("lock expired")
	LockCancel  = MutexErr("lock cancel")
)

type opt struct {
	mxer rsync.Lockable
	key  string
	d    time.Duration
}

type Mutex struct {
	opt
}

type MutexOptions func(op *opt)

func NewMutex(opts ...MutexOptions) *Mutex {
	option := &opt{}
	for _, op := range opts {
		op(option)
	}

	return &Mutex{opt: *option}
}

func (m *Mutex) Lock(ctx context.Context) {
	for {
		ok := m.TryLock(ctx)
		if ok {
			return
		}
	}
}

func (m *Mutex) Unlock(ctx context.Context) {
	m.mxer.Del(ctx, m.key)
}

func (m *Mutex) TryLock(ctx context.Context) bool {
	now := time.Now().Unix()

	return m.mxer.SetNX(ctx, m.key, now, m.d)
}
