package main

import (
	"fmt"
	"os"

	"rtty/baudot"
)

func main() {
	baudRate := baudot.BAUD_DELAY_45

	// Print named file, if one is provided. Otherwise, print test
	b, err := baudot.New(baudRate)
	if err != nil {
		fmt.Printf("%s: unable to create a New() context: %v\n", os.Args[0], err)
		os.Exit(1)
	}

	// TODO: Add command line parsing to deal with options, like speed,
	// print test, etc...
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
	const (
		spaces   = "SpaceTest:B         1         2         3         4         5    E"
		puncts   = "Punctuation: !@#$%^&*()_+={}[]:;\"',.?/  !@#$%^&*()_+={}[]:;\"',.?/  !@#$%^&*()_+={}[]:;\"',.?/ "
		figsLtrs = "Z0Y1X2W3V4U5T6S7R8P9Q0Z0Y1X2W3V4U5T6S7R8P9Q0Z0Y1X2W3V4U5T6S7R8P9Q0"
		test     = "the quick brown fox jumped over the lazy dog's back     1234567890\n" +
			"ryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryry\n" +
			"sgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsg"
	)

	c.PrintLine(spaces)
	c.PrintLine(puncts)
	c.PrintLine(figsLtrs)
	c.PrintLine("")
	c.PrintLine(test)
	c.PrintLine(test)
	c.PrintLine(test)
	c.PrintLine(test)
	c.PrintLine("\n\n")
}
