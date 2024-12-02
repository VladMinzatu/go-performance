package main

import (
	"fmt"
	"go-performance/utils"
	"sync"
)

func main() {
	numGoroutines := 10_000   // we'll average over a larger number of goroutines for numerical stability
	var block <-chan struct{} // nil channel to block all goroutines while we take the memory measurement

	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	memBefore := utils.GetMemoryStats().Sys

	launchGoroutines(numGoroutines, &wg, block) // launch numGoroutines goroutines that are empty and blocked
	wg.Wait()                                   // wait for all goroutines to be created

	memAfter := utils.GetMemoryStats().Sys
	fmt.Printf("%.3f kb memory used per goroutine\n", float64(memAfter-memBefore)/float64(numGoroutines)/1000)
}

func launchGoroutines(numGoroutines int, wg *sync.WaitGroup, block <-chan struct{}) {
	noop := func() { // these are our goroutines - empty and doing nothing, to measure the goroutine memory footprint
		wg.Done()
		<-block // wait here
	}

	for i := 0; i < numGoroutines; i++ {
		go noop()
	}
}
