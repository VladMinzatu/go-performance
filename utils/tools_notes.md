_Disclaimer: Feel free to read this doc, but these are just notes I took for myself as I worked through these examples. Don't expect complete and clear info_

# Tools for profiling and observability

## The `runtime` package

Package runtime contains operations that interact with Go's runtime, including functions that control goroutines, the GC and allow access to low level metrics such as memory usage (see pkg.go.dev/runtime)

## Benchmarks

Benchmarks execute a certain code segment a number of times in order to get a stable estimate for the execution time of that code segment.

The easiest and standard way to write benchmarks in Go is to use the built-in functionality in the `testing` package, i.e. writing a specific type of test (inside a `*_test.go` file) and then executing it with a `go test -bench=.` command.

Benchmark functions have a signature of the form `func Benchmark*(*testing.B)`. `b.N` will typically be used inside the function and determines the number of iterations. The output from running `go test -bench=.` will tell us how many times the code segment under test was run and what was the average time it took.

## Profilers

Go supports multiple profilers listed below. They way to work with them typically involves these general steps:

- modify your code to start a profiler (this will typically involve pointing it to a file that collects the profile data)
- run the code
- use the `go tool pprof -http:8080 <filename>` to inspect the output

This can be combined nicely with benchmarks when appropriate (e.g. `go test -bench='.' -cpuprofile='cpu.prof' -memprofile='mem.prof'`), but sometimes you'll probably want to test large code segments or an entire application in a production environment or something that resembles that.

Note: `go test ...` would also work with those \*profile options.

### CPU profiler

Will capture statistics about the on-CPU time of the code by interrupting the program every 10ms through interrupts and taking a stack trace.

The profiling is started by placing this in the function that includes the code you want to profile:

```
...
pprof.StartCPUProfile(f)
defer pprof.StopCPUProfile()
```

or with https://github.com/pkg/profile :

```
defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
```

### Memory profiler

Collects data about the memory allocations per stack trace. The go runtime itself does this by recording the stack trace that lead to allocations (at a certain sample rate, which is tunable, capturing all is possible by setting `runtime.MemProfileRate`), instead of using an OS interrupt to capture CPU cycles.

Example using https://github.com/pkg/profile

```
defer profile.Start(profile.MemProfile, profile.MemProfileRate(1), profile.ProfilePath(".")).Stop()
```

Alternatively, `pprof.WriteHeapProfile(f)` could be used to do a heap dump at a point in time. `go tool pprof mem.prof` can be used on the output.

### Block profiler

Will capture off-CPU time spent waiting on channels and mutexes (but not sleep, I/O, GC). The statistics will show cumulatie delays per stack trace.

### Mutex profiler

The same as the block profiler, but only looks at mutexes, but excluding channels.

Subtle difference, though. The mutex profiler compiles statistics per stack trace _causing the blocking_, as opposed to being blocked (which is what the block profiler does). Generally, you'd want to use both profilers together.

### Goroutine profiler

Collects data about number of goroutines per stack trace.

WARNING: this is a stop the world profiler, no sampling mechanism.

But it can be useful for detecting goroutine leaks or diagnose why a program might hang.

## Tracing: Go's built-in runtime Scheduler Tracer

Tracing is the recording of timestamped events. (this is useful at the go application level, the same way distributed tracing is useful for understanding performance at a distributed system level)

The built in runtime tracer captures scheduler, GC, contention, syscal etc. events.(see src/runtime/trace.go)

This can expose issues not captured by the profilers mentioned above.

To run, e.g. using https://github.com/pkg/profile :

```
defer profile.Start(profile.TraceProfile, profile.Path(".")).Stop()
```

It comes with its own UI to view thread and goroutine timelines.
After running, `go tool trace trace.out`

Can also be exported to Prometheus as metrics accessible through the metrics endpoint (e.g. https://github.com/MadhavJivrajani/gse - essentially it runs a go program with `GODEBUG=schedtrace=10 <binary>` and then it scans the stderr for "SCHED" and then parses those traces to extract metrics and pushes them to prometheus).


## Observability

For higher level observability, check https://github.com/open-telemetry/opentelemetry-go 

# Ideas

- Optimizing memory allocations (heap vs stack) - benchmarks

- stack allocations escaping to the heap

- Memory leak scenarios

- Goroutine leaks

- performance channel sync vs locking primitives (channel sync is seamless as part of the regular scheduling, whereas locking is extra and expensive operatio - demo it)
- goroutine vs mutex use cases

## References

https://blog.logrocket.com/benchmarking-golang-improve-function-performance/

https://www.youtube.com/watch?v=7hg4T2Qqowk

https://www.youtube.com/watch?v=nok0aYiGiYA

https://stackademic.com/blog/profiling-go-applications-in-the-right-way-with-examples
