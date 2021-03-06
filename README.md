# gctx
[![Build Status](https://github.com/ansiwen/gctx/workflows/CI/badge.svg?branch=master)](https://github.com/ansiwen/gctx/actions?query=branch%3Amaster)
[![Go Report Card](https://goreportcard.com/badge/github.com/ansiwen/gctx)](https://goreportcard.com/report/github.com/ansiwen/gctx)
[![GoDoc](https://pkg.go.dev/badge/github.com/ansiwen/gctx?status.svg)](https://pkg.go.dev/github.com/ansiwen/gctx?tab=doc)

Implementation of a goroutine context (`gctx`).

> **TL;DR** \
> Set a goroutine context with `gctx.Set(ctx)` and retrieve it from anywhere in
> the same or a child goroutine with `gctx.Get()`.

A gctx is a `context.Context` that is neither global, nor in the local function
scope, but a property of a specific goroutine. It can only be created and
retrieved by code being executed within the according goroutine. The gctx can be
set with `gctx.Set()`:

```golang
ctx := context.WithValue(context.Background(), myCtxVal{}, "foobar")
go func() {
    gctx.Set(ctx)
    ...
    /* some code that eventually calls doSomething() */
    ...
}()
```

Then, from any function called within the same or a child goroutine this gctx
can be retrieved with `gctx.Get()`:

```golang
func doSomething() {
    ...
    ctx := gctx.Get()
    v := ctx.Value(myCtxVal{})
    log.Printf("ctx: %v", v) // prints "ctx: foobar"
    ...
}
```


## Why
The Go language decided for good reasons that goroutines are [anonymous by
design](https://go.dev/doc/faq#no_goroutine_id). So there is no ID or any
goroutine-local storage exposed. (Well, almost, see below.) It is advised to
always use a `context.Context` as the first parameter in all functions that
require knowledge about a specific context.

However, there are also good reasons to have such a goroutine specific context
generally available. The most common reason is logging:

- There are libraries that don't support a context parameter (yet). When such a
  library is producing logs, one might want to annotate the log message with the
  context in which the library call has been made. But without the possibility
  to "tunnel" a context through the library to the log callback, that is not
  possible.
- Debug logs or tracing appear *anywhere* in your code. Annotating these logs
  with the current context would require to add a context parameter to *every*
  *single* *function*, which is simply unnecessary boilerplate and in most cases
  not feasible.

`gctx` allows to solve these situations by attaching a context to the current
goroutine, which can be retrieved from any code, that is running in the same
goroutine or a child, with a simple `gctx.Get()` call. It is effectively a
goroutine-local storage, and could also be used to identify a goroutine, which
is not recommended. Therefore it should only be used, if there is really no
other possibility available.

## How
In fact, goroutines *do* [have an
ID](https://github.com/golang/go/blob/851ecea4cc99ab276109493477b2c7e30c253ea8/src/runtime/runtime2.go#L438),
but it is hidden and not accessible with any exported API. That's why there have
been a couple of attempts to circumvent that, by implementing a goroutine-local
storage in order to solve above problems, like for example:
 - https://github.com/jtolio/gls
 - https://github.com/tylerb/gls
 - https://github.com/modern-go/gls

But all these packages are pretty inefficient and/or are highly dependent on
implementation details.

Fortunately, there is another goroutine property that *is* - at least partly -
exported: [profiling labels](https://pkg.go.dev/runtime/pprof). They can be
stored with
[`runtime/pprof.SetGoroutineLabels()`](https://pkg.go.dev/runtime/pprof#SetGoroutineLabels)
in the goroutine-local storage and are automatically inherited by child
goroutines. But again, there is no exported API that allows to read these label
from the goroutine local storage.

So, in order to make (ab)use of these labels for a `gctx`, it is required to do
one little runtime hack: the internal functions
[`setProfLabel()`](https://github.com/golang/go/blob/851ecea4cc99ab276109493477b2c7e30c253ea8/src/runtime/proflabel.go#L12)
and
[`getProfLabel()`](https://github.com/golang/go/blob/851ecea4cc99ab276109493477b2c7e30c253ea8/src/runtime/proflabel.go#L38)
must be made accessible with a `//go:linkname` instruction. (Of course this is
also an implementation detail, that could break anytime. But since it is part of
an exported functionality, it can't just disappear, and a fix should always be
easily possible.)

These functions are managing a pointer to a `labelMap` (which is a
`map[string]string`) in the goroutine-local storage. `gctx` extends this type to
piggyback a `context.Context` on it:

```golang
struct {
	labelMap
	context.Context
}
```

It stores pointers to objects of this extended yet compatible type instead of
the original `labelMap`. A special label in the `labelMap` is set that
indicates, that a context is available. That way the attached context is - like
the labels - automatically memory managed and inherited by child goroutines. 

## Usage

See
[documentation](https://pkg.go.dev/github.com/ansiwen/gctx#section-documentation)
and [examples](https://pkg.go.dev/github.com/ansiwen/gctx#pkg-examples).
