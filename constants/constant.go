package main

import "fmt"

const LogFile = "/var/log/log.txt" //Untyped String

const LogFileStr string = "/var/log/log_str.txt" //Typed String

const ageLimit = 30 //untyped int

const total = iota //untyped int

const (
	Pi = 3.14 //untyped float
)

const ( //1KB & 1MB which is represented in bytes binary - Bitshifting
	KB = 1 << 10
	MB = KB << 10
	GB = MB << 10
)

const (
	fileSizeLimitBytes     = 5 * MB                 //5242880 - Multiplication
	halfFileSizeLimitBytes = fileSizeLimitBytes / 2 //2621440 - Division
)

type Week int

const (
	Monday    Week = iota //0
	Tuesday               //1
	Wednesday             //2
	Thursday              //3
	Friday                //4
	Saturday              //5
	Sunday                //6
)

func main() {

	desc("LogFile", LogFile)
	desc("ageLimit", ageLimit)
	desc("Pi", Pi)

	desc("Thursday", Thursday)

	desc("total", total)

	desc("KB", KB)
	desc("MB", MB)
	desc("GB", GB)
}

func desc(varName string, v interface{}) {
	fmt.Printf("Var : %s | Type : %T | Value : %v\n", varName, v, v)
}
