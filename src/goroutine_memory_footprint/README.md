# The memory footprint of a Goroutined

Goroutines are advertised as being very lightweight, both in terms of memory usage and context switching. In this script, we validate this claim by lauching many minimal goroutines and checking their average memory footprint.

## How to run it

Simply run the code file in this directory:

```
go run main.go
```

It should print out the average memory usage of an empty goroutine. On my system, it's roughly 2.5 kb, which is pretty lightweight indeed!

## References

A similar experiment is done in the excellent book [Concurrency in Go: Tools and Techniques for Developers](https://isbnsearch.org/isbn/1491941197) by Katherine Cox-Buday, which is where the inspiration for this test comes from.
