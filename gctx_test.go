package gctx

import (
	"context"
	"fmt"
	"runtime/pprof"
	"strconv"
	"sync"
	"testing"
	"unsafe"
)

func assert(t *testing.T) func(bool, ...interface{}) {
	return func(ok bool, args ...interface{}) {
		t.Helper()
		if !ok {
			t.Error(args...)
		}
	}
}

func assertB(b *testing.B) func(bool, ...interface{}) {
	return func(ok bool, args ...interface{}) {
		b.Helper()
		if !ok {
			b.Error(args...)
		}
	}
}

type ctxKey struct{} // nolint

func TestGctxSet(t *testing.T) {
	assert := assert(t)
	ctx := context.WithValue(context.Background(), ctxKey{}, "test")
	ctx2 := Get()
	assert(ctx2 == nil)
	Set(ctx)
	ctx2 = Get()
	assert(ctx2 != nil)
	v := ctx2.Value(ctxKey{})
	assert(v == "test")
}

func TestGctxSetWithLabel(t *testing.T) {
	assert := assert(t)
	ctx := context.WithValue(context.Background(), ctxKey{}, "test")
	ctx = pprof.WithLabels(ctx, pprof.Labels("foo", "bar"))
	Set(ctx)
	lm := (*labelMap)(runtimeGetProfLabel())
	assert(lm != nil)
	assert(len(*lm) == 2)
	_, ok := (*lm)[labelTag]
	assert(ok)
	assert((*lm)["foo"] == "bar")
	gctx := Get()
	assert(gctx != nil)
	v := gctx.Value(ctxKey{})
	assert(v == "test")
	v, _ = pprof.Label(gctx, "foo")
	assert(v == "bar")
}

func TestGctxSetLabelID(t *testing.T) {
	assert := assert(t)
	ctx := context.WithValue(context.Background(), ctxKey{}, "test")
	Set(ctx)
	lm := (*labelMap)(runtimeGetProfLabel())
	assert(lm != nil)
	id, ok := strconv.ParseUint((*lm)[labelTag], 16, int(unsafe.Sizeof(lastID)))
	assert(ok == nil)
	assert(id > 0)
	Set(ctx)
	lm = (*labelMap)(runtimeGetProfLabel())
	assert(lm != nil)
	id2, ok := strconv.ParseUint((*lm)[labelTag], 16, int(unsafe.Sizeof(lastID)))
	assert(ok == nil)
	assert(id2 != id)
}

func TestGctxDoubleSet(t *testing.T) {
	assert := assert(t)
	ctx := context.WithValue(context.Background(), ctxKey{}, "test")
	ctx2 := context.WithValue(ctx, ctxKey{}, "test2")
	Set(ctx)
	gctx := Get()
	assert(gctx != nil)
	v := gctx.Value(ctxKey{})
	assert(v == "test")
	Set(ctx2)
	gctx = Get()
	assert(gctx != nil)
	v = gctx.Value(ctxKey{})
	assert(v == "test2")
}

func TestGctxMulti(t *testing.T) {
	assert := assert(t)
	var wg, wg2, wg3 sync.WaitGroup
	do := func(s string) {
		ctx := context.WithValue(context.Background(), ctxKey{}, s)
		wg.Add(1)
		wg3.Add(1)
		go func() {
			Set(ctx)
			wg.Done()
			wg2.Wait()
			gctx := Get()
			assert(gctx != nil)
			v := gctx.Value(ctxKey{})
			assert(v == s)
			wg3.Done()
		}()
	}
	wg2.Add(1)
	for i := 1; i < 10; i++ {
		do(fmt.Sprintf("test%d", i))
	}
	wg.Wait()
	wg2.Done()
	wg3.Wait()
}

func BenchmarkGctx(b *testing.B) {
	assert := assertB(b)
	var wg sync.WaitGroup
	do := func(s string) {
		ctx := context.WithValue(context.Background(), ctxKey{}, s)
		wg.Add(1)
		go func() {
			Set(ctx)
			gctx := Get()
			assert(gctx != nil)
			v := gctx.Value(ctxKey{})
			assert(v == s)
			wg.Done()
		}()
	}
	for i := 0; i < b.N; i++ {
		do(fmt.Sprintf("test%d", i))
	}
	wg.Wait()
}

func TestGctxInherit(t *testing.T) {
	assert := assert(t)
	var wg, wg2, wg3 sync.WaitGroup
	ctx := context.WithValue(context.Background(), ctxKey{}, "test")
	wg.Add(1)
	wg2.Add(1)
	wg3.Add(1)
	go func() {
		Set(ctx)
		ctx := Get()
		assert(ctx != nil)
		v := ctx.Value(ctxKey{})
		assert(v == "test")
		go func() {
			wg2.Wait()
			ctx := Get()
			assert(ctx != nil)
			v := ctx.Value(ctxKey{})
			assert(v == "test")
			wg3.Done()
		}()
		wg.Done()
	}()
	wg.Wait()
	wg2.Done()
	wg3.Wait()
}
