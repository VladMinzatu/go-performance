In Go versions prior to 1.23, one source of memory leaks was the use of timer.After() inside a for(-select) loop.
The issue comes from the fact that such a timer, once launched, consumes some heap memory until the timer expires.

Go 1.23 changed that, and now the memory used by the timer is available to be reclaimed when the timer is no longer referenced by any variables still in scope.

The program inside this directory puts this to the test with a minimalistic reproduction of the issue, without the use of extra goroutines and a for-select.

## How to run it

Simply run the code file in this directory:

```
go run main.go
```

On my machine, this prints out:

```
Before:
119.136 kb allocated
After:
124.872 kb allocated
```

That is a modest difference, indeed. But let's compare with older versions of Go. I tried running this same code using Go 1.21.1 and got the following output:

```
Before:
124.968 kb allocated
After:
21079.336 kb allocated
```

So that checks out and we can safely use time.After() without worry of runaway memory usage beyond what we can see from regular scoping rules.

If you have to use a version of Go prior to 1.23, the common fix for this issue is to use time.NewTimer() instead. The pattern looks like this:

```
duration := time.Hour
timer := time.NewTimer(duration)
for {
  timer.Reset(duration)
  select {
    ...
    case <-timer.C:
    ...
  }
}
```

Using this pattern both in Go 1.23, as well as in older versions incurs minimal memory overhead:

```
Before:
119.232 kb allocated
After:
119.656 kb allocated
```

But as I said, as of Go 1.23, this is no longer necessary, and using time.After() should be slightly more convenient.

## References

I first read about this subtle issue in the excellent book [100 Go Mistakes and How to Avoid Them](https://isbnsearch.org/isbn/1617299596) by Teiva Harsanyi.
