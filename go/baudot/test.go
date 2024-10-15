package baudot

import (
	"fmt"
)

func (c *Convert) PrintTest() {
	const (
		spaces   = "SpaceTest:B         1         2         3         4         5    E"
		puncts   = "Punctuation: !@#$%^&*()_+={}[]:;\"',.?/  !@#$%^&*()_+={}[]:;\"',.?/  !@#$%^&*()_+={}[]:;\"',.?/ "
		figsLtrs = "Z0Y1X2W3V4U5T6S7R8P9Q0Z0Y1X2W3V4U5T6S7R8P9Q0Z0Y1X2W3V4U5T6S7R8P9Q0"
		test     = "the quick brown fox jumped over the lazy dog's back     1234567890\n" +
			"ryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryry\n" +
			"sgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsg"
	)

	fmt.Fprintf(c, "%s\n", spaces)
	fmt.Fprintf(c, "%s\n", puncts)
	fmt.Fprintf(c, "%s\n", figsLtrs)
	fmt.Fprintf(c, "%s\n", "")
	fmt.Fprintf(c, "%s\n", test)
	fmt.Fprintf(c, "%s\n", test)
	fmt.Fprintf(c, "%s\n", test)
	fmt.Fprintf(c, "%s\n", test)
	fmt.Fprintf(c, "%s\n", "\n\n")
}
