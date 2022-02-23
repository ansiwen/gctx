// Package gctx implements a goroutine context (gctx). A gctx is a non-global
// context that is linked to one specific goroutine and can be retrieved from
// any scope within that goroutine.
package gctx

import (
	"context"
	"runtime/pprof"
	"strconv"
	"sync/atomic"
	"unsafe"
)

const labelTag = "github.com/ansiwen/gctx"

var lastID uintptr

type labelMapAndContext struct {
	lm  labelMap
	ctx context.Context
}

// Set ctx as the current goroutine's gctx. Goroutines created within this
// goroutine inherit the gctx.
func Set(ctx context.Context) {
	var gctx labelMapAndContext
	pprof.SetGoroutineLabels(ctx)
	lm := (*labelMap)(runtimeGetProfLabel())
	if lm != nil && len(*lm) > 0 {
		gctx.lm = *lm
	} else {
		gctx.lm = make(labelMap)
	}
	id := atomic.AddUintptr(&lastID, 1)
	gctx.lm[labelTag] = strconv.FormatUint(uint64(id), 16)
	gctx.ctx = ctx
	runtimeSetProfLabel(unsafe.Pointer(&gctx))
}

// Get the current goroutine's gctx, which has been set in the same or a parent
// goroutine with gctx.Set() before.
func Get() context.Context {
	lm := (*labelMap)(runtimeGetProfLabel())
	if lm != nil {
		if _, ok := (*lm)[labelTag]; ok {
			gctx := (*labelMapAndContext)(unsafe.Pointer(lm))
			return gctx.ctx
		}
	}
	return nil
}
