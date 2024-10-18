package main

import (
	"fmt"
	"os"

	"rtty/baudot"
)

func main() {
	baudRate := baudot.BaudDelay45

	file := ""
	if len(os.Args) > 1 {
		// Print named file
		_, err := os.Stat(os.Args[1])
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		file = os.Args[1]
	}

	b, err := baudot.New(baudRate)
	if err != nil {
		fmt.Printf("%s: unable to create a New() context: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	defer b.Close()

	if file != "" {
		if err := b.PrintFile(os.Args[1]); err != nil {
			fmt.Printf("Cannot print file %q: %v\n", os.Args[1], err)
			os.Exit(1)
		}
	} else {
		// print test data
		b.PrintTest()
	}

	// TODO: Add command line parsing to deal with options, like speed,
	// print test, etc...
}
