package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var size = 50000000

func main() {

	debug.SetGCPercent(-1) //stop auto gc

	_ = make([]int, size)
	printMemUsage("Using Int")

	_ = make([]float64, size)
	printMemUsage("Using Float")

	_ = make([]complex128, size)
	printMemUsage("Using Complex")

	_ = make([]bool, size)
	printMemUsage("Using Bool")

	_ = make([]byte, size)
	printMemUsage("Using Byte")

	_ = make([]struct{}, size)
	printMemUsage("Using Empty Structs")
}

func printMemUsage(v string) {

	defer runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Print(v + fmt.Sprintf(" - Memory Allocated : %v MB\n", m.Alloc/1024/1024))

}
