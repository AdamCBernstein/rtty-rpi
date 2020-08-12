package main

import (
	"rtty/baudot"
)

func main() {
	spaces := "SpaceTest:                                                      \n"
	I0 := "Z0Y1X2W3V4U5T6S7R8P9Q0Z0Y1X2W3V4U5T6S7R8P9Q0Z0Y1X2W3V4U5T6S7R8P9Q0\n"
	test := "the quick brown fox jumped over the lazy dog's back     1234567890\n" +
		  "ryryryryryryryryryryryryryryryryryryryryrryyryryryryryryryryryryry\n" +
		  "sgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsg\n"

	baudot := baudot.New()
	baudot.Print(spaces)
	baudot.Print(I0)
	baudot.Print(test)
	baudot.Print(test)
	baudot.Print(test)
	baudot.Print(test)
	baudot.Print("\n\n")
}
