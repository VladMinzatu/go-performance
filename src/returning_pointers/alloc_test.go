package main

import (
	"testing"
)

var resultData byte // benchmarks assign values to this var in order to avoid optimisations due to unused results
const dataSize = 100

type ResultData struct {
	data [dataSize]byte
}

type ResultSlice []byte

func BenchmarkReturnDataValue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		result := getDataArray()
		resultData = result.data[0]
	}
}

func BenchmarkReturnDataPointer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		result := getDataArrayPointer()
		resultData = result.data[0]
	}
}

func BenchmarkReturnSlice(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		data := dataSlice(dataSize)
		resultData = data[0]
	}
}

func getDataArray() ResultData {
	var data ResultData
	return data
}

func getDataArrayPointer() *ResultData {
	var data ResultData
	return &data
}

func dataSlice(n int) []byte {
	var data = make([]byte, n)
	return data
}
