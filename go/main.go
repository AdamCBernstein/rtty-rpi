package main

import (
	"rtty/baudot"
)

func main() {
	test := "the quick brown fox jumped over the lazy dog's back 1234567890 times\n" +
		"ryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryry"

	baudot := baudot.New()
	baudot.Print(test)
	baudot.Print(test)
	baudot.Print(test)
	baudot.Print(test)
	baudot.Print(test)
	baudot.Print("\n\n")
}
