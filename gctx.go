// Package gctx implements a goroutine context (gctx). A gctx is a non-global
// context that is linked to one specific goroutine and can be retrieved from
// any scope within that goroutine.
package gctx

import (
	"context"
	"runtime/pprof"
)

// Do calls f with ctx set as the current goroutine's gctx. Goroutines created
// within f inhert the gctx, but it is only available as long as f didn't
// return. If this is not sufficient, because the child outlives f, the child
// has to create its own gctx, which can be based on the gctx of f.
func Do(ctx context.Context, f func()) {
	oldLabels := runtimeGetProfLabel()
	id := newID()
	gctx := pprof.WithLabels(ctx, label(id))
	add(gctx, id)
	pprof.SetGoroutineLabels(gctx)
	f()
	runtimeSetProfLabel(oldLabels)
	delete(id)
}

// Get retrieves the current goroutines gctx. It returns nil if no gctx has been
// created, or it has already been deleted.
func Get() context.Context {
	return get(currentID())
}
