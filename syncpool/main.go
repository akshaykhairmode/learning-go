package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

type User struct {
	Name    []byte
	Age     int
	Address []byte
	Bio     []byte
	Phone   []byte
}

var pool = &sync.Pool{
	New: func() interface{} {
		return &User{
			Name:    make([]byte, 64<<11),
			Age:     20,
			Address: make([]byte, 64<<11),
			Bio:     make([]byte, 64<<11),
			Phone:   make([]byte, 64<<11),
		}
	},
}

var f, err = os.OpenFile("log", os.O_WRONLY, 0644)

func main() {

	if err != nil {
		panic(err)
	}

}

func newUser() *User {
	return &User{
		Name:    make([]byte, 64<<11),
		Age:     20,
		Address: make([]byte, 64<<11),
		Bio:     make([]byte, 64<<11),
		Phone:   make([]byte, 64<<11),
	}
}

func userCreator(c int) {
	wg := &sync.WaitGroup{}
	for i := 0; i < c; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			user := newUser()
			_ = user
		}()
	}
	wg.Wait()
}

func userCreatorPool(c int) {
	wg := &sync.WaitGroup{}
	for i := 0; i < c; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			user := pool.Get().(*User)
			pool.Put(user)
		}()
	}
	wg.Wait()
}

func PrintMemUsage(s string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Fprintf(f, s+" :: ")
	fmt.Fprintf(f, "Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Fprintf(f, "\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Fprintf(f, "\tSys = %v MiB", bToMb(m.Sys))
	fmt.Fprintf(f, "\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
