package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {

	c := make(chan string, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			c <- strconv.Itoa(i)
		}(i)
	}

	go func() {
		fmt.Println("sleeping for 1 second")
		time.Sleep(1 * time.Second)
		close(c)
		fmt.Println("channel closed")
	}()

	for v := range c {
		fmt.Println(v)
	}

}
