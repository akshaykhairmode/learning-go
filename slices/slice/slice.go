package main

import "fmt"

func main() {

	slice := []int{1, 2, 3, 4, 5}
	desc(slice)

	otherSlice := slice[0:2]
	desc(otherSlice)

	slice[0] = 100
	desc(slice)
	desc(otherSlice)

	otherSlice[0] = 50
	desc(slice)

	// output
	// Type is []int, Address is 0xc000020090, Capacity : 5, Length : 5, value is [1 2 3 4 5]
	// Type is []int, Address is 0xc000020090, Capacity : 5, Length : 2, value is [1 2]
	// Type is []int, Address is 0xc000020090, Capacity : 5, Length : 5, value is [100 2 3 4 5]
	// Type is []int, Address is 0xc000020090, Capacity : 5, Length : 2, value is [100 2]
	// Type is []int, Address is 0xc000020090, Capacity : 5, Length : 5, value is [50 2 3 4 5]

}

func desc(t []int) {
	fmt.Printf("Type is %T, Address is %p, Capacity : %v, Length : %v, value is %v\n", t, t, cap(t), len(t), t)
}
