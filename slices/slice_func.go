package main

import (
	"fmt"
)

var cnt = 1

func main() {

	slice := []string{"London", "Washington", "New York"}

	desc(slice, "slice")
	updateIndex(slice)
	desc(slice, "slice")

	otherSlice := slice[0:2]
	desc(otherSlice, "otherSlice")
	updateIndex(otherSlice)
	otherSlice = append(otherSlice, "Mumbai")
	desc(otherSlice, "otherSlice")
	otherSlice = append(otherSlice, "Bangalore")
	desc(otherSlice, "otherSlice")
	otherSlice = append(otherSlice, "Hyderabad")
	desc(otherSlice, "otherSlice")
	desc(slice, "slice")

	// output
	// 1. Printing slice, Address is 0xc000070150, Capacity : 3, Length : 3, value is [London Washington New York]
	// 2. Printing inside updateIndex, Address is 0xc00005a180, Capacity : 6, Length : 6, value is [Bermuda Triangle Washington New York One Two Three]
	// 3. Printing slice, Address is 0xc000070150, Capacity : 3, Length : 3, value is [Bermuda Triangle Washington New York]
	// 4. Printing otherSlice, Address is 0xc000070150, Capacity : 3, Length : 2, value is [Bermuda Triangle Washington]
	// 5. Printing inside updateIndex, Address is 0xc00005a1e0, Capacity : 6, Length : 5, value is [Bermuda Triangle Washington One Two Three]
	// 6. Printing otherSlice, Address is 0xc000070150, Capacity : 3, Length : 3, value is [Bermuda Triangle Washington Mumbai]
	// 7. Printing otherSlice, Address is 0xc00005a240, Capacity : 6, Length : 4, value is [Bermuda Triangle Washington Mumbai Bangalore]
	// 8. Printing otherSlice, Address is 0xc00005a240, Capacity : 6, Length : 5, value is [Bermuda Triangle Washington Mumbai Bangalore Hyderabad]
	// 9. Printing slice, Address is 0xc000070150, Capacity : 3, Length : 3, value is [Bermuda Triangle Washington Mumbai]
}

func updateIndex(v []string) {
	v[0] = "Bermuda Triangle"
	v = append(v, "One", "Two", "Three")
	desc(v, "inside updateIndex")
}

func desc(t []string, sl string) {
	fmt.Printf("%d. Printing %v, Address is %p, Capacity : %v, Length : %v, value is %v\n", cnt, sl, t, cap(t), len(t), t)
	cnt++
}
