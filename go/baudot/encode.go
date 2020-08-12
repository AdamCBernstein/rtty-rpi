package baudot

import (
	"fmt"
	"unicode"

	"rtty/gpio"
)

type baudotBits [8]bool
type Convert struct {
	shift bool
}

func shiftLettersFigures(shift bool, c *Convert, retBits []baudotBits) []baudotBits {
	if shift {
		return append(retBits, baudotConv[shiftUp])
	} else {
		return append(retBits, baudotConv[shiftDown])
	}
}

func asciiToBaudot(r rune, c *Convert) ([]baudotBits, bool) {
	var retBits = make([]baudotBits, 0, 8)

	// Deal with control characters first
	switch r {
	case lineFeed, carriageReturn:
		c.shift = false
		return append(retBits,
			baudotConv[carriageReturn],
			baudotConv[lineFeed],
			baudotConv[shiftDown]), true

	case spaceChar:
		return append(retBits, baudotConv[spaceChar]), true
	}

	// Get Baudot bits value for ASCII character
	bits, ok := baudotConv[r]
	if !ok {
		return nil, false
	}

	if shift := bits[LTRS_FIGS_BIT]; shift != c.shift {
		c.shift = shift
		retBits = shiftLettersFigures(c.shift, c, retBits)
	}
	return append(retBits, bits), true
}

func printRune(r rune, c *Convert) {
	rUpper := unicode.ToUpper(r)
	bitsSlice, ok := asciiToBaudot(rUpper, c)
	if !ok {
		return
	}

	// Give some console feedback. TODO: Make callback function
	fmt.Printf("%c", rUpper)
	writeBits(bitsSlice)
}

func writeBits(bitsChar []baudotBits) {
	for _, bits := range bitsChar {
		for _, bit := range bits[:LTRS_FIGS_BIT] {
			wiringpi.WriteBit(bit)
			wiringpi.DelayMicroseconds(BAUD_DELAY_45)
		}
		wiringpi.DelayMicroseconds((BAUD_DELAY_45 * 50) / 100)
	}
}

func initializeTeletype(c *Convert) {
	printRune(shiftDown, c)
	printRune(shiftDown, c)
	printRune(carriageReturn, c)
	printRune(lineFeed, c)
	printRune(shiftDown, c)
}
