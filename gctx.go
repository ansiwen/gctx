// Package gctx implements a goroutine context (gctx). A gctx is a non-global
// context that is linked to one specific goroutine and can be retrieved from
// any scope within that goroutine.
package gctx

import (
	"context"
	"runtime/pprof"
	"unsafe"
)

const labelTag = "github.com/ansiwen/gctx"

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
	gctx.lm[labelTag] = ""
	gctx.ctx = ctx
	runtimeSetProfLabel(unsafe.Pointer(&gctx))
}

// Get retrieves the current goroutines gctx.
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
