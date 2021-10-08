package main

import (
	"fmt"
	"strconv"
)

type employee struct {
	id         int
	firstName  string
	lastName   string
	middleName string
	age        int
}

func (this employee) String() string {
	return "EMP-" + strconv.Itoa(this.id) + " (" + this.firstName + " " + this.middleName + " " + this.lastName + ")"
}

func main() {

	emp := employee{id: 1001, firstName: "Tom", middleName: "Marvolo", lastName: "Riddle"}
	print(emp)
	printAddr(&emp)

	fmt.Println(fmt.Sprintf("Sprintf : Your subscription is ended : %v", emp))
}

func print(v interface{}) {

	fmt.Printf("Type : %T\n", v)

	fmt.Println("Println :", v)

	fmt.Printf("Printf %%v : %v\n", v)
	fmt.Printf("Printf %%s : %s\n", v)

	fmt.Printf("Printf %%+v : %+v\n", v)
	fmt.Printf("Printf %%#v : %#v\n", v)
}

func printAddr(v interface{}) {
	fmt.Printf("Printf Address %%v : %v\n", v)
	fmt.Printf("Printf Address Ptr %%p : %p\n", v)
}
