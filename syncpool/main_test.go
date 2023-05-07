package main

import "testing"

var num = 100

func BenchmarkWithoutPool(b *testing.B) {
	PrintMemUsage("Before without pool")
	for i := 0; i < b.N; i++ {
		userCreator(num)
	}
	PrintMemUsage("After without pool")
}

func BenchmarkWithPool(b *testing.B) {
	PrintMemUsage("Before with pool")
	for i := 0; i < b.N; i++ {
		userCreatorPool(num)
	}
	PrintMemUsage("After with pool")
}
