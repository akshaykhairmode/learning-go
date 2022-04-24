package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

//our custom type
type Number interface {
	int | uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64 | float32 | float64
}

func main() {

	//INT
	fmt.Println("SumWithConstraints", SumWithConstraints(5, 5, 5, 5))
	fmt.Println("SumWithSpecificTypes", SumWithSpecificTypes(5, 5, 5, 5))
	fmt.Println("SumWithCustomType", SumWithCustomType(5, 5, 5, 5))

	//Float
	fmt.Println("SumWithConstraints", SumWithConstraints(5.5, 22.22, 33.33, 44.44))
	fmt.Println("SumWithSpecificTypes", SumWithSpecificTypes(5.5, 22.22, 33.33, 44.44))
	fmt.Println("SumWithCustomType", SumWithCustomType(5.5, 22.22, 33.33, 44.44))

}

//If we use constraints package.
func SumWithConstraints[T constraints.Integer | constraints.Float](i ...T) T {
	var r T
	for _, v := range i {
		r += v
	}
	return r
}

//If we do not use constraints package
func SumWithSpecificTypes[T int | uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64 | float32 | float64](i ...T) T {
	var r T
	for _, v := range i {
		r += v
	}
	return r
}

//If we define a custom type
func SumWithCustomType[T Number](i ...T) T {
	var r T
	for _, v := range i {
		r += v
	}
	return r
}
