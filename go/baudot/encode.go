package baudot

import (
	wp "rtty/gpio"
	"unicode"
)

type baudotBits [8]bool
type convert struct {
	shift bool
}

func asciiToBaudot(r rune, c *convert) ([]baudotBits, bool) {
	var retBits = make([]baudotBits, 0, 8)

	// Deal with control characters first
	switch r {
	case lineFeed, carriageReturn:
		retBits = append(retBits, baudotConv[carriageReturn])
		retBits = append(retBits, baudotConv[lineFeed])
		retBits = append(retBits, baudotConv[shiftDown])
		c.shift = false
		return retBits, true
	case spaceCharacter:
		retBits = append(retBits, baudotConv[spaceCharacter])
		return retBits, true
	}

	// Get Baudot bits value for ASCII character
	bits, ok := baudotConv[r]
	if !ok {
		return nil, false
	}

	shift := bits[LTRS_FIGS_BIT]
	if shift != c.shift {
		c.shift = shift
		if shift {
			retBits = append(retBits, baudotConv[shiftUp])
		} else {
			retBits = append(retBits, baudotConv[shiftDown])
		}
	}
	return append(retBits, bits), true
}

func printRune(r rune, c *convert) {
	bitsSlice, ok := asciiToBaudot(unicode.ToUpper(r), c)
	if !ok {
		return
	}
	writeBits(bitsSlice)
}

func writeBits(bitsChar []baudotBits) {
	for _, bits := range bitsChar {
		for _, bit := range bits[:LTRS_FIGS_BIT] {
			wp.WriteBit(bit)
			wp.DelayMicroseconds(BAUD_DELAY_45)
		}
		wp.DelayMicroseconds((BAUD_DELAY_45 * 50) / 100)
	}
}

func initializeTeletype(c *convert) {
	printRune(shiftDown, c)
	printRune(shiftDown, c)
	printRune(carriageReturn, c)
	printRune(lineFeed, c)
	printRune(shiftDown, c)
}
