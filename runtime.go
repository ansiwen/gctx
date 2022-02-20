package gctx

import (
	_ "unsafe" // for go:linkname
)

//go:linkname runtimeGetProfLabel runtime/pprof.runtime_getProfLabel
func runtimeGetProfLabel() *map[string]string
