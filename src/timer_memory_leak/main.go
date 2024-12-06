package main

import (
	"fmt"
	"go-performance/utils"
	"runtime"
	"time"
)

func main() {
	numTimers := 100_000

	fmt.Println("Before:")
	printRuntimeStats(utils.GetMemoryStats())

	launchTimers(numTimers)

	fmt.Println("After:")
	printRuntimeStats(utils.GetMemoryStats())
}

func printRuntimeStats(memStats runtime.MemStats) {
	formattedMemoryUsage := fmt.Sprintf("%.3f kb", float64(memStats.Alloc)/1000)
	fmt.Println(formattedMemoryUsage, "allocated")
}

func launchTimers(numTimers int) {
	for i := 0; i < numTimers; i++ {
		time.After(time.Hour)
	}
}
