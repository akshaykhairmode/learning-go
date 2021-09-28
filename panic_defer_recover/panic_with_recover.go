package main

import "fmt"

func main() {

	fmt.Println("Before funcOne")
	funcOne()
	fmt.Println("After funcOne")
}

func funcOne() {

	fmt.Println("Before funcTwo")
	funcTwo()
	fmt.Println("After funcTwo")
}

func funcTwo() {

	fmt.Println("Before funcThree")
	funcThree()
	fmt.Println("After funcThree")

}

func funcThree() {
	defer myRecover()

	panic("Panic Inside funcThree")
}

func myRecover() {
	if err := recover(); err != nil {
		fmt.Println("Recovered Error :", err)
	}
}
