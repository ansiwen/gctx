package gctx

import (
	"runtime/pprof"
	"strconv"
	"sync/atomic"
	"unsafe"
)

const labelKey = "gctx_id"

var lastID uintptr

func newID() uintptr {
	return atomic.AddUintptr(&lastID, 1)
}

func currentID() uintptr {
	lbls := runtimeGetProfLabel()
	if lbls != nil {
		if v, ok := (*lbls)[labelKey]; ok {
			id, _ := strconv.ParseUint(v, 16, int(unsafe.Sizeof(lastID)))
			return uintptr(id)
		}
	}
	return 0
}

func label(id uintptr) pprof.LabelSet {
	v := strconv.FormatUint(uint64(id), 16)
	return pprof.Labels(labelKey, v)
}
