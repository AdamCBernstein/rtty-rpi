package main

import (
	"fmt"
	"os"

	"rtty/baudot"
)

func main() {
	// Print named file, if one is provided. Otherwise, print test
	b := baudot.New()
	if len(os.Args) > 1 {
		if err := b.PrintFile(os.Args[1]); err != nil {
			fmt.Printf("Cannot print file %q: %v\n", os.Args[1], err)
			os.Exit(1)
		}
	} else {
		printTest(b)
	}
}

func printTest(c *baudot.Convert) {
	spaces := "SpaceTest:B         1         2         3         4         5    E"
	I0 := "Z0Y1X2W3V4U5T6S7R8P9Q0Z0Y1X2W3V4U5T6S7R8P9Q0Z0Y1X2W3V4U5T6S7R8P9Q0"
	test := "the quick brown fox jumped over the lazy dog's back     1234567890\n" +
		"ryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryry\n" +
		"sgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsg"

	c.PrintLine(spaces)
	c.PrintLine(I0)
	c.PrintLine("")
	c.PrintLine(test)
	c.PrintLine(test)
	c.PrintLine(test)
	c.PrintLine(test)
	c.PrintLine("\n\n")
}
