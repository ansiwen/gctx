package gctx_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/ansiwen/gctx"
	"github.com/ansiwen/gctx/internal/utils"
)

func TestGctx(t *testing.T) {
	assert := utils.Assert(t)
	var wg sync.WaitGroup
	wg.Add(2)
	ctx1 := context.WithValue(context.Background(), utils.CtxKey{}, "test1")
	ctx2 := context.WithValue(context.Background(), utils.CtxKey{}, "test2")
	go gctx.Do(ctx1, func() {
		ctx := gctx.Get()
		assert(ctx != nil)
		v := ctx.Value(utils.CtxKey{})
		assert(v == "test1")
		wg.Done()
	})
	go gctx.Do(ctx2, func() {
		var wg2 sync.WaitGroup
		wg2.Add(1)
		go func() {
			ctx := gctx.Get()
			assert(ctx != nil)
			v := ctx.Value(utils.CtxKey{})
			assert(v == "test2")
			wg2.Done()
		}()
		ctx := gctx.Get()
		assert(ctx != nil)
		v := ctx.Value(utils.CtxKey{})
		assert(v == "test2")
		wg2.Wait()
		wg.Done()
	})
	wg.Wait()
}

type writer func(p []byte) (n int, err error)

func (w writer) Write(p []byte) (n int, err error) {
	return w(p)
}

func Example_logger() {
	type requestID struct{}
	// "global" logger instance
	log := func(s string) {
		rid, _ := gctx.Get().Value(requestID{}).(int)
		fmt.Printf("[LOG] reqID: %d - %s\n", rid, s)
	}
	// placeholder for an external function without context support doing an
	// arbitrary request and using a global logger
	createSomeResource := func() {
		log("sent create request")
	}
	var wg sync.WaitGroup
	// asynchronously request 5 times to create some resource
	for i := 0; i < 5; i++ {
		wg.Add(1)
		ctx := context.WithValue(context.Background(), requestID{}, i)
		go gctx.Do(ctx, func() {
			createSomeResource()
			wg.Done()
		})
	}
	wg.Wait()
	// Unordered output:
	// [LOG] reqID: 0 - sent create request
	// [LOG] reqID: 1 - sent create request
	// [LOG] reqID: 2 - sent create request
	// [LOG] reqID: 3 - sent create request
	// [LOG] reqID: 4 - sent create request
}
