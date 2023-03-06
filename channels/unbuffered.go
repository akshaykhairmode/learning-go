package main

import (
	"log"
	"time"
)

func main() {

	c := make(chan string, 1)

	go func() {
		for data := range c {
			log.Println(data)
		}
	}()

	c <- "test"

	time.Sleep(1 * time.Second)

}
