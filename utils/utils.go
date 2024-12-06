package utils

import (
	"runtime"
)

func GetMemoryStats() runtime.MemStats {
	runtime.GC()
	var s runtime.MemStats
	runtime.ReadMemStats(&s)
	return s
}
