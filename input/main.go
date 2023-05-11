package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func main() {
	readPassword()
}

func readInput() {
	// print our prompt message
	fmt.Print("Enter your name : ")

	// Start a reader, we can use scanner also in case we want to keep reading multiple times.
	r := bufio.NewReader(os.Stdin)

	// read input till user presses enter, readstring will include the newline as well.
	inp, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}

	// strip the new line
	inp = strings.TrimRight(inp, "\n")

	// If name is empty show below message
	if inp == "" {
		fmt.Println("Welcome : Stranger")
		return
	}

	fmt.Printf("Welcome : %s\n", inp)
}

func readPassword() {
	fmt.Println("Enter password:")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}

	if string(password) == "abilityrush" {
		fmt.Println("Authenticated")
		return
	}

	fmt.Println("Invalid password")
}
