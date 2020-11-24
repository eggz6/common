package util

import (
	"context"
)

// SafeGo logf is called when it is retrieved, so make it simple without any panic
func SafeGo(ctx context.Context, do func(), logf func(string, ...interface{})) {
	if do == nil {
		return
	}

	go func(ctx context.Context) {
		defer func() {
			if e := recover(); e != nil {
				if logf != nil {
					logf("safeg recover. ctx=%v, msg=%v", ctx, e)
				}
			}
		}()
		do()
	}(ctx)
}
