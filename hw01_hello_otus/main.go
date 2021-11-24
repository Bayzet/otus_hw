package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

const msg string = "Hello, OTUS!"

func main() {
	reversed := stringutil.Reverse(msg)
	fmt.Println(reversed)
}
