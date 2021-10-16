package main

import "fmt"

type PhoneModel int

const (
	iPhone4  PhoneModel = iota //0
	iPhone5                    //1
	iPhone6                    //2
	iPhone7                    //3
	iPhone8                    //4
	iPhone10                   //5
	iPhone11                   //6
	iPhone12                   //7
	iPhone13                   //8
)

// func (pm PhoneModel) String() string {

// 	switch pm {
// 	case 0:
// 		return "iPhone4"
// 	case 1:
// 		return "iPhone5"
// 	case 2:
// 		return "iPhone6"
// 	case 3:
// 		return "iPhone7"
// 	case 4:
// 		return "iPhone8"
// 	case 5:
// 		return "iPhone9"
// 	case 6:
// 		return "iPhone10"
// 	case 7:
// 		return "iPhone11"
// 	case 8:
// 		return "iPhone12"
// 	case 9:
// 		return "iPhone13"
// 	}

// 	return ""
// }

func main() {
	fmt.Printf("Some Error occurred when listing phone : %v , error : %v\n", iPhone11, "Some error")
}
