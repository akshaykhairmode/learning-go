package main

import "fmt"

func main() {

	fmt.Println("Before funcOne")
	funcOne()
	fmt.Println("After funcOne")
}

func funcOne() {

	defer myRecover()

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
	panic([]int{1, 2, 3})
}

func myRecover() {
	if err := recover(); err != nil {
		fmt.Println("Recovered Error :", err)
	}
}
