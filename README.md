# gctx

Implementation of a goroutine context (`gctx`).

A gctx is a `context.Context` that is neither global, nor in the local function
scope, but a property of a specific goroutine. It can only be created and
retrieved by code being executed within the according goroutine.

## Why
The Go language decided for good reasons that goroutines are [anonymous by
design](https://go.dev/doc/faq#no_goroutine_id). So there is no ID or any
goroutine-local storage exposed. (Well, almost, see below.) It is advised to
always use a `context.Context` as the first parameter in all functions that
require knowledge about a goroutine specific context.

However, there are also good reasons to have such a goroutine specific context
generally available. The most common one is logging: For example there are many
libraries that don't support a context parameter (yet), but they allow to hook
your own code into it by providing interfaces or callbacks. If such code is not
strictly serialized (which is the usual case in concurrent code), there is no
possibility to correlate the callbacks with the originating call into the
library, without the possibility to "tunnel" a context through the 3rd party
library. Also in your own code, logging, especially debug logs or tracing, can
appear everywhere in your code. That means to realize context aware debug
logging, a context must be added to *every* *single* function, which is simply
unnecessary boilerplate and in most cases not feasible.

`gctx` allows to solve these situations by attaching a context to the current
goroutine, which can be retrieved from any code that is running in the
goroutine. It is effectively a goroutine local storage. Therefore it is advised
to only use it, if there is no other possibility available.

## How
There have been a couple of attempts to implement goroutine-local storage in
order to solve these problems, like:
 - https://github.com/jtolio/gls
 - https://github.com/tylerb/gls
 - https://github.com/modern-go/gls

 but all these packages are pretty inefficient and/or are highly dependent on
 implementation details. In fact, goroutines *do* [have an
ID](https://github.com/golang/go/blob/851ecea4cc99ab276109493477b2c7e30c253ea8/src/runtime/runtime2.go#L438),
but it is hidden and not accessible with any exported API. But there is another
goroutine property that *is* - at least partly - exported: [profiling
labels](https://pkg.go.dev/runtime/pprof). They can be stored with
[`runtime/pprof.SetGoroutineLabels()`](https://pkg.go.dev/runtime/pprof#SetGoroutineLabels)
in the goroutine local storage and therefore can be used to store an identifier
for a goroutine specific context. However, there is no exported API that allows
to read these label from the goroutine local storage. In order to make (ab)use
of these labels for a `gctx`, it is required to do one little hack: the internal
function
[`getProfLabel()`](https://github.com/golang/go/blob/851ecea4cc99ab276109493477b2c7e30c253ea8/src/runtime/proflabel.go#L38)
must be made accessible with a single `//go:linkname` instruction. Of course
this is also an implementation detail, than can break anytime. But since it is
part of an exported functionality, it can't just disappear, and a fix will
always be easily possible.

## Usage

See examples and tests.
