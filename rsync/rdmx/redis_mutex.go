package rdmx

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type RdMutex struct {
	rd redis.Cmdable
	cd chan int
}

func NewRdMutex(rd redis.Cmdable) *RdMutex {
	cd := make(chan int, 1)

	return &RdMutex{rd: rd, cd: cd}
}

func (rx *RdMutex) SetNX(ctx context.Context, key string, val interface{}, d time.Duration) bool {
	select {
	case rx.cd <- 1:
		ok, err := rx.rd.SetNX(key, val, d).Result()
		if err != nil {
			return false
		}

		return ok
	default:
		return false
	}
}

func (rx *RdMutex) Del(ctx context.Context, key string) {
	_, _ = rx.rd.Del(key).Result()
}
