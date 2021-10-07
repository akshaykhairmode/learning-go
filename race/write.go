package main

import (
	"fmt"
	"sync"
)

func main() {

	fmt.Println("Start")

	m := make(map[int]struct{})

	wg := &sync.WaitGroup{}
	mx := &sync.Mutex{}

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go updateMap(wg, mx, m, i)
	}

	wg.Wait()

	fmt.Println(m)
}

func updateMap(wg *sync.WaitGroup, mx *sync.Mutex, m map[int]struct{}, r int) {
	defer wg.Done()
	mx.Lock()
	defer mx.Unlock()
	m[r] = struct{}{}
}
