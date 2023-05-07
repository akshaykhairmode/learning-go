package main

import (
	"fmt"

	"github.com/akshaykhairmode/learning-go/build/internal"
)

var ver string

func main() {
	fmt.Println(ver)
	fmt.Println(internal.BuildTime)
}
