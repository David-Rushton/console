package main

import (
	"fmt"
	"os"
	"strings"

	"example.com/console/console/markup"
)

var input string

func init() {
	if len(os.Args) < 2 {
		fmt.Println("Usage")
		fmt.Println()
		fmt.Println(`console "text to print"`)
		fmt.Println()

		os.Exit(0)
	}

	input = strings.Join(os.Args[1:], input)
}

func main() {

	// console.Clear()
	output := markup.Parse(input)
	fmt.Printf("%v\n", output)
}
