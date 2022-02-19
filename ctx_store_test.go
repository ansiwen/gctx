package gctx

import (
	"context"
	"testing"

	"github.com/ansiwen/gctx/internal/utils"
)

func TestStore(t *testing.T) {
	assert := utils.Assert(t)
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.CtxKey{}, "test")
	add(ctx, 42)
	ctx = get(42)
	assert(ctx != nil)
	v := ctx.Value(utils.CtxKey{})
	assert(v == "test")
	ctx = get(43)
	assert(ctx == nil)
	delete(42)
	ctx = get(42)
	assert(ctx == nil)
}
