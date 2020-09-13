package rdmx

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/elliotchance/redismock"
	"github.com/go-redis/redis"
)

func Test_SetNX(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	rd := redismock.NewNiceMock(client)
	rd.On("SetNX", "test", 1, time.Duration(0)).Return(redis.NewBoolResult(true, nil))
	rx := NewRdMutex(rd)

	ch := rx.SetNX(context.TODO(), "test", 1, 0)
	res := <-ch
	if !res {
		t.Fatal("set failed")
	}
}

func Test_SetNXCancel(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	rd := redismock.NewNiceMock(client)
	rd.On("SetNX", "test", 1, time.Duration(0)).Return(redis.NewBoolResult(true, nil))
	rx := NewRdMutex(rd)

	ctx, cancel := context.WithCancel(context.TODO())
	ch := rx.SetNX(ctx, "test", 1, 0)
	cancel()
	_, ok := <-ch
	if ok {
		t.Fatal("set failed")
	}
}
