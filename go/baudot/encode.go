package baudot

type baudotBits [8]bool

func (c *Convert) asciiToBaudot(r rune) ([]baudotBits, bool) {
	retBits := make([]baudotBits, 0, 3)

	// Deal with control characters first
	switch r {
	case carriageReturn, lineFeed:
		// Force downshift on CR/LF to keep Teletype shift state in-sync
		c.shift = false
		return append(retBits,
			baudotConv[carriageReturn],
			baudotConv[lineFeed],
			baudotConv[shiftDown]), true
	case spaceChar:
		return []baudotBits{baudotConv[spaceChar]}, true
	}

	// Get Baudot bits value for ASCII character
	if bits, ok := baudotConv[r]; !ok {
		return nil, false
	} else {
		if bitsLtrsFigs, ok := c.shiftLettersFigures(retBits, bits[LTRS_FIGS_BIT]); !ok {
			return []baudotBits{bits}, true
		} else {
			return append(bitsLtrsFigs, bits), true
		}
	}
}

func (c *Convert) shiftLettersFigures(retBits []baudotBits, shift bool) ([]baudotBits, bool) {
	if c.shift == shift {
		return nil, false
	}

	c.shift = shift
	if shift {
		return append(retBits, baudotConv[shiftUp]), true
	} else {
		return append(retBits, baudotConv[shiftDown]), true
	}
}
