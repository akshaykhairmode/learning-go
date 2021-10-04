package main

import (
	"bufio"
	"os"
	"strconv"
)

func main() {

}

func writeToFile() {

	file := "writeSlow.txt"
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer os.Remove(file)
	defer f.Close()

	for i := 0; i < 100000; i++ {
		f.WriteString(strconv.Itoa(i))
	}

}

func writeToFileBuffered() {

	file := "writeFast.txt"
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer os.Remove(file)
	defer f.Close()

	buf := bufio.NewWriter(f)

	for i := 0; i < 100000; i++ {
		buf.WriteString(strconv.Itoa(i))
	}

	buf.Flush()

}
