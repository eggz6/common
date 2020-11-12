package poker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func Test_NewRes(t *testing.T) {
	res := newres()

	if len(res.ch) != 0 {
		t.Fatal("new res failed")
	}
}

func Test_TryExecuteMul(t *testing.T) {
	res := newres()

	table, _ := newTable(3)
	err := res.decorator(context.Background(), newGame("1", table), nil)

	if err != InvalidResHandle {
		t.Fatalf("try execute mul failed. should be err=%v", InvalidResHandle)
	}

	err = res.decorator(context.Background(), newGame("1", table),
		func(ctx context.Context, g *game) error {
			return nil
		})

	if err != nil {
		t.Fatalf("try execute mul failed")
	}

	err = res.decorator(context.Background(), newGame("1", table),
		func(ctx context.Context, g *game) error {
			return nil
		})

	if err != HasBeenExecuted {
		t.Fatalf("try execute mul failed. should be err=%v", HasBeenExecuted)
	}
}

func Test_P_TryExecute(t *testing.T) {
	res := newres()
	count := uint32(0)

	var wg sync.WaitGroup

	t.Run("group", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			t.Run(fmt.Sprintf("test_%v", i), func(tt *testing.T) {
				tt.Parallel()
				ok := res.tryExecute()
				if !ok {
					atomic.AddUint32(&count, 1)
				}

				tt.Logf("parallel ok=%v", ok)
				wg.Done()
			})
		}
	})

	wg.Wait()

	if count != 9 {
		t.Fatalf("tryExecute failed. count=%v", count)
	}
}
