package gctx

import (
	"context"
	"runtime/pprof"
	"testing"

	"github.com/ansiwen/gctx/internal/utils"
)

func TestCurrentID(t *testing.T) {
	assert := utils.Assert(t)
	ctx := context.Background()
	l := pprof.Labels("foo", "bar")
	ctx = pprof.WithLabels(ctx, l)
	ctx = context.WithValue(ctx, utils.CtxKey{}, "foobar")
	pprof.SetGoroutineLabels(ctx)
	id := currentID()
	assert(id == 0)
	ctx = pprof.WithLabels(ctx, label(42))
	pprof.SetGoroutineLabels(ctx)
	id = currentID()
	assert(id == 42)
	done := make(chan struct{})
	go func() {
		id := currentID()
		assert(id == 42)
		close(done)
	}()
	<-done
	id = currentID()
	assert(id == 42)
	l = pprof.Labels(labelKey, "brokenGctxId")
	pprof.SetGoroutineLabels(pprof.WithLabels(ctx, l))
	id = currentID()
	assert(id == 0)
}

func TestNewID(t *testing.T) {
	assert := utils.Assert(t)
	id := newID()
	assert(id != 0)
	id2 := newID()
	assert(id != id2)
	assert(id2 != 0)
}
