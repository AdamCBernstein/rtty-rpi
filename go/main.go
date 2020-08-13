package main

import (
	"os"

	"rtty/baudot"
)

func main() {
	// Print named file, if one is provided. Otherwise, print test
	b := baudot.New()
	if len(os.Args) > 1 {
		b.PrintFile(os.Args[1])
	} else {
		b.PrintTest()
	}
}
