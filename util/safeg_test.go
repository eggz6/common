package util

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func Test_Safeg(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	SafeGo(context.TODO(), func() {
		fmt.Println("goroutine here")
		panic("panic here")
	}, func(f string, args ...interface{}) {
		wg.Done()
		t.Logf(f, args...)
	})

	wg.Wait()
}
