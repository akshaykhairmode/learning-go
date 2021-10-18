package main

import (
	"encoding/json"
	"fmt"
)

type shop struct {
	Name             string `json:"shopName"`
	Address          string `json:"shopAddr"`
	Typ              string `json:"type"`
	Phone            int    `json:"-"`
	AverageCustomers *int   `json:"avgCust"`
}

func main() {

	//Without Fields, will follow fields as per order of definition
	shop1 := shop{"Shop1", "Canada", "Grocery", 1111111111, new(int)}

	//With fields
	shop2 := shop{Name: "Shop2", Address: "US", Typ: "Pharmacy", Phone: 1111111111, AverageCustomers: new(int)}

	//New will return a pointer to shop, deferencing is done by compiler in case of structs
	shop3 := new(shop)
	shop3.Name = "Shop3"
	shop3.Address = "Japan"
	shop3.Typ = "Restaurant"
	shop3.Phone = 1111111111
	shop3.AverageCustomers = new(int)

	fmt.Println("Printing Initial Values...")
	printShop(shop1, shop2, *shop3)
	fmt.Println("Updating Customer Values...")
	updateCustomers(shop1, 5)
	updateCustomers(shop2, 10)
	updateCustomers(*shop3, 15)
	fmt.Println("Printing Values After Customer Update...")
	printShop(shop1, shop2, *shop3)
	fmt.Println("Updating Phone Values...")
	updatePhone(shop1, 000000000)
	updatePhone(shop2, 000000000)
	updatePhone(*shop3, 000000000)
	printShop(shop1, shop2, *shop3)

	var js []byte
	js, _ = json.Marshal(shop1)
	fmt.Printf("Json Marshal Shop1 : %s\n", js)

	js, _ = json.Marshal(shop2)
	fmt.Printf("Json Marshal Shop2 : %s\n", js)

	js, _ = json.Marshal(shop3)
	fmt.Printf("Json Marshal Shop3 : %s\n", js)

}

//Wont update as we are updating a copy of shop
func updatePhone(s shop, v int) {
	s.Phone = v
}

//Will update as the copy is referencing to an actual location
func updateCustomers(s shop, v int) {
	*s.AverageCustomers = v
}

func printShop(s ...shop) {

	for _, v := range s {
		fmt.Printf("Printing for : %v ---- %+v", v.Name, v)
		fmt.Printf(" - AverageCustomers Value - %v \n", *v.AverageCustomers)
	}

}
