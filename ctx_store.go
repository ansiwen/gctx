package gctx

import (
	"context"
	"sync"
)

var ctxStore sync.Map

func add(ctx context.Context, id uintptr) {
	ctxStore.Store(id, ctx)
}

func get(id uintptr) context.Context {
	v, ok := ctxStore.Load(id)
	if !ok {
		return nil
	}
	return v.(context.Context)
}

func delete(id uintptr) {
	ctxStore.Delete(id)
}
