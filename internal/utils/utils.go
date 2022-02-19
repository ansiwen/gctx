// Package utils ...
package utils

import "testing"

// Assert ...
func Assert(t *testing.T) func(bool, ...interface{}) {
	return func(ok bool, args ...interface{}) {
		t.Helper()
		if !ok {
			t.Error(args...)
		}
	}
}

// CtxKey ...
type CtxKey struct{} // nolint
