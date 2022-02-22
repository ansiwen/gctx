package gctx_test

import (
	"context"
	"fmt"
	"sync"

	"github.com/ansiwen/gctx"
)

func Example_logger() {
	type requestID struct{}

	// "global" logger instance
	log := func(s string) {
		ctx := gctx.Get()
		reqID, _ := ctx.Value(requestID{}).(int)
		fmt.Printf("[LOG] reqID: %d - %s\n", reqID, s)
	}

	// placeholder for an external library function without context support
	// doing something and using a global logger
	someLibraryFunction := func() {
		/* ... */
		log("did something")
		/* ... */
	}

	var wg sync.WaitGroup
	// asynchronously call a library function for 5 times
	for i := 0; i < 5; i++ {
		wg.Add(1)
		ctx := context.WithValue(context.Background(), requestID{}, i)
		go func() {
			gctx.Set(ctx)
			someLibraryFunction()
			wg.Done()
		}()
	}
	wg.Wait()

	// Unordered output:
	// [LOG] reqID: 0 - did something
	// [LOG] reqID: 1 - did something
	// [LOG] reqID: 2 - did something
	// [LOG] reqID: 3 - did something
	// [LOG] reqID: 4 - did something
}
