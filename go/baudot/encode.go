package baudot

type baudotBits [8]bool

func (c *Convert) shiftLettersFigures(retBits []baudotBits, shift bool) []baudotBits {
	if c.shift == shift {
		return nil
	}

	c.shift = shift
	if shift {
		return append(retBits, baudotConv[shiftUp])
	} else {
		return append(retBits, baudotConv[shiftDown])
	}
}

func (c *Convert) asciiToBaudot(r rune) ([]baudotBits, bool) {
	var retBits = make([]baudotBits, 0, 8)

	// Deal with control characters first
	switch r {
	case carriageReturn, lineFeed:
		// Force downshift on CR/LF to keep Teletype shift state in-sync
		return append(retBits,
			baudotConv[carriageReturn],
			baudotConv[lineFeed],
			baudotConv[shiftDown]), true
	case spaceChar:
		return append(retBits, baudotConv[spaceChar]), true
	}

	// Get Baudot bits value for ASCII character
	if bits, ok := baudotConv[r]; !ok {
		return nil, false
	} else {
		return append(c.shiftLettersFigures(retBits, bits[LTRS_FIGS_BIT]), bits), true
	}
}
