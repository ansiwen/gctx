// Package gctx implements a goroutine context.
package gctx

import (
	"context"
	"runtime/pprof"
)

// Do calls f with ctx set as the current goroutine context.
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

// Get retrieves the current goroutine context.
func Get() context.Context {
	return get(currentID())
}
