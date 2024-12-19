It is a common optimization tip in Go to avoid returning pointers from a function, unless there is a good reason for it that doesn't have to do with optimization.
It is said that this typically holds true for data structures up to 10 MB, which is just about all data structures you'll typically use.

But as a handy corrolary, if you're ever in a special case where you're passing around multi-MB data structures, just favor passing around pointers always.

This observation is made because passing pointers *into* functions is often regarded as more performant than passing the value itself, because it avoids creating a copy of the value.

The reason why returning a pointer may perform worse than returning a value is because when a pointer is returned from a function, the variable that is pointed to escapes to the heap (the compiler cannot guarantee that it won't be accessed in an unsafe way if it were to be allocated on the stack, so it will be allocated on the heap). In addition to the extra work created for the GC, allocating on the heap itself is a much more time consuming operation than allocating on the stack (where we just have a pointer handy to the next available memory address).

Note that this is not an issue when sharing down (i.e. passing a pointer into a function), only when sharing up. Because when the pointer is shared up, it is pointing to a no longer valid stack area.

Slice variables may not escape to the heap if their size is known to the compiler. But if the size is based on a variable, or if the backing array is reallocated following an `append`, it will be allocated ont he heap.

## How to run it

```
go test  -gcflags=-l  -bench=.
```

Note the usage of the `-gcflags=-l` flag which disables all inlining. Alternatively, we could use the pragma `//go:noinline` above each of the functions returning 

If neither of these techniques are used, the compiler will optimize these simple demonstrative cases by inlining the functions that produce data and the results won't be reproducible.

The output on my machine looks like this:
```
BenchmarkReturnDataValue-10             216829573                5.386 ns/op           0 B/op          0 allocs/op
BenchmarkReturnDataPointer-10           49921564                22.72 ns/op          112 B/op          1 allocs/op
BenchmarkReturnSlice-10                 51316957                22.59 ns/op          112 B/op          1 allocs/op
PASS
ok      go-performance/src/returning_pointers   4.734s
```
We can see that returning the value (size 100 bytes in this case) is around 4 times faster than returning the pointer. And we can also see the extra heap allocation in the case of returning a pointer, which causes the increased latency. We can also see that when we return a slice of unknown size at compile time, the heap allocation is not avoided.

We can also try rerunning the benchmarks without the  `-gcflags=-l` flag to see what the inlining optimization does for us in normal cases:
```
BenchmarkReturnDataValue-10             1000000000               0.3163 ns/op          0 B/op          0 allocs/op
BenchmarkReturnDataPointer-10           1000000000               0.3184 ns/op          0 B/op          0 allocs/op
BenchmarkReturnSlice-10                 52528479                22.34 ns/op          112 B/op          1 allocs/op
PASS
ok      go-performance/src/returning_pointers   2.259s
```
Nothing changes in the case of the slice, but the inlining of the function gives us more than an order of magnitude performance boost by also avoiding the data copying when returning a value. Thanks, Go compiler!

## References

This topic is discussed in the following two books (along with others and many online resources):
- [Learning Go: An Idiomatic Approach to Real-World Go Programming](https://isbnsearch.org/isbn/1098139291)
- [100 Go Mistakes and How to Avoid Them](https://isbnsearch.org/isbn/1617299596) 