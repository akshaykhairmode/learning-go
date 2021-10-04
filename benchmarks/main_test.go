package main

import "testing"

func Benchmark_writeToFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		writeToFile()
	}
}

func Benchmark_writeToFileBuffered(b *testing.B) {
	for i := 0; i < b.N; i++ {
		writeToFileBuffered()
	}
}
