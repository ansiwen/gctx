package gctx_test

import (
	"context"
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
