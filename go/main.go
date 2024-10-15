package main

import (
	"fmt"
	"os"

	"rtty/baudot"
)

func main() {
	baudRate := baudot.BaudDelay45

	// Print named file, if one is provided. Otherwise, print test
	b, err := baudot.New(baudRate)
	if err != nil {
		fmt.Printf("%s: unable to create a New() context: %v\n", os.Args[0], err)
		os.Exit(1)
	}
	defer b.Close()

	// TODO: Add command line parsing to deal with options, like speed,
	// print test, etc...
	if len(os.Args) > 1 {
		if err := b.PrintFile(os.Args[1]); err != nil {
			fmt.Printf("Cannot print file %q: %v\n", os.Args[1], err)
			os.Exit(1)
		}
	} else {
		b.PrintTest()
	}
}
