package baudot

func (c *Convert) PrintTest() {
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
