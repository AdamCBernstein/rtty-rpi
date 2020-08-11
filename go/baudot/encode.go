package baudot

import (
	"fmt"
	wp "rtty/gpio"
	"unicode"
)

type Baudot interface {
	Print()
}

type baudotBits [8]bool

type convert struct {
	baudotTable []baudotBits
	shift       bool
}

const I = true
const O = false
const shiftUp = '\x02'
const shiftDown = '\x01'
const lineFeed = '\n'
const carriageReturn = '\r'
const spaceCharacter = ' '
const nullChar = '\xfe'

// ITA-2 / US-TTY table. Reference:
// https://en.wikipedia.org/wiki/Baudot_code#ITA_2_and_US-TTY
//
// format: bit0 = up/down shift
//         bit1-bit7 = start/data5/stop
var baudotConv = map[rune]baudotBits{
	// All entries here are LETTERS; bit0 = O
	/* 0   */ 'A': {O, O, I, I, O, O, O, I}, /* A / - */
	/* 1   */ 'B': {O, O, I, O, O, I, I, I}, /* B / ? */
	/* 2   */ 'C': {O, O, O, I, I, I, O, I}, /* C / : */
	/* 3   */ 'D': {O, O, I, O, O, I, O, I}, /* D / $ */
	/* 4   */ 'E': {O, O, I, O, O, O, O, I}, /* E / 3 */
	/* 5   */ 'F': {O, O, I, O, I, I, O, I}, /* F */
	/* 6   */ 'G': {O, O, O, I, O, I, I, I}, /* G */
	/* 7   */ 'H': {O, O, O, O, I, O, I, I}, /* H */
	/* 8   */ 'I': {O, O, O, I, I, O, O, I}, /* I / 8 */
	/* 9   */ 'J': {O, O, I, I, O, I, O, I}, /* J / ' */
	/*10   */ 'K': {O, O, I, I, I, I, O, I}, /* K / ( */
	/*11   */ 'L': {O, O, O, I, O, O, I, I}, /* L / ) */
	/*12   */ 'M': {O, O, O, O, I, I, I, I}, /* M / . */
	/*13   */ 'N': {O, O, O, O, I, I, O, I}, /* N / , */
	/*14   */ 'O': {O, O, O, O, O, I, I, I}, /* O / 9 */
	/*15   */ 'P': {O, O, O, I, I, O, I, I}, /* P / O */
	/*16   */ 'Q': {O, O, I, I, I, O, I, I}, /* Q / 1  */
	/*17   */ 'R': {O, O, O, I, O, I, O, I}, /* R / 4 */
	/*18   */ 'S': {O, O, I, O, I, O, O, I}, /* S */
	/*19   */ 'T': {O, O, O, O, O, O, I, I}, /* T / 5 */
	/*20   */ 'U': {O, O, I, I, I, O, O, I}, /* U / 7 */
	/*21   */ 'V': {O, O, O, I, I, I, I, I}, /* V / ; */
	/*22   */ 'W': {O, O, I, I, O, O, I, I}, /* W / 2 */
	/*23   */ 'X': {O, O, I, O, I, I, I, I}, /* X / / */
	/*24   */ 'Y': {O, O, I, O, I, O, I, I}, /* Y / 6 */
	/*25   */ 'Z': {O, O, I, O, O, O, I, I}, /* Z / " */
	/*26*/ '\xfe': {O, O, O, O, O, O, I, I}, /* NULL */
	/*27  */ '\n': {O, O, I, O, O, O, I, I}, /* LF/ LF */
	/*28   */ ' ': {O, O, O, I, O, O, I, I}, /* SPACE / SPACE */
	/*29  */ '\r': {O, O, O, O, I, O, I, I}, /* CR / CR */
	/*30 */ '\x02': {O, I, I, O, I, I, I, I}, /* SHIFT_UP */
	/*31 */ '\x01': {O, I, I, I, I, I, I, I}, /* SHIFT_DOWN */
	/*32 */ '\x00': {O, O, O, O, O, O, O, O}, /* Open */
	/*33 */ '\xff': {I, I, I, I, I, I, I, I}, /* closed */

	// All entries here are FIGS; bit0 = 1
	/* 0*/ '-': {I, O, I, I, O, O, O, I}, /*  A / - */
	/* 1*/ '?': {I, O, I, O, O, I, I, I}, /*  B / ? */
	/* 2*/ ':': {I, O, O, I, I, I, O, I}, /*  C / : */
	/* 3*/ '$': {I, O, I, O, O, I, O, I}, /*  D / $ */
	/* 4*/ '3': {I, O, I, O, O, O, O, I}, /*  E / 3 */
	/* 5*/ // BEL or '
	/* 6*/ '\a': {I, O, O, I, O, I, I, I}, /* G / BEL */
	/* 7*/ // unknown
	/* 8*/ '8': {I, O, O, I, I, O, O, I}, /*  I  / 8 */
	/* 9*/ '\'': {I, O, I, I, O, I, O, I}, /* J  / ' */
	/*10*/ '(': {I, O, I, I, I, I, O, I}, /*  K / ( */
	/*11*/ ')': {I, O, O, I, O, O, I, I}, /*  L / ) */
	/*12*/ '.': {I, O, O, O, I, I, I, I}, /*  M / . */
	/*13*/ ',': {I, O, O, O, I, I, O, I}, /*  N / , */
	/*14*/ '9': {I, O, O, O, O, I, I, I}, /*  O / 9 */
	/*15*/ '0': {I, O, O, I, I, O, I, I}, /*  P / O */
	/*16*/ '1': {I, O, I, I, I, O, I, I}, /*  Q / 1  */
	/*17*/ '4': {I, O, O, I, O, I, O, I}, /*  R / 4 */
	/*19*/ '5': {I, O, O, O, O, O, I, I}, /*  T / 5 */
	/*20*/ '7': {I, O, I, I, I, O, O, I}, /*  U / 7 */
	/*21*/ ';': {I, O, O, I, I, I, I, I}, /*  V / ; */
	/*22*/ '2': {I, O, I, I, O, O, I, I}, /*  W / 2 */
	/*23*/ '/': {I, O, I, O, I, I, I, I}, /*  X / / */
	/*24*/ '6': {I, O, I, O, I, O, I, I}, /*  Y / 6 */
	/*25*/ '"': {I, O, I, O, O, O, I, I}, /*  Z / " */
	// These are duplicates for FIGURES
	// /*26*/	'\xff': {I, O, O, O, O, O, I, I}, /* NULL */
	// /*27*/  '\n': {I, O, I, O, O, O, I, I}, /* LF/ LF */
	// /*28*/	' ': {I, O, O, I, O, O, I, I}, /* SPACE / SPACE */
	// /*29*/	'\r': {I, O, O, O, I, O, I, I}, /* CR / CR */
	// /*30*/	'\x01': {I, I, I, O, I, I, I, I}, /* SHIFT_UP */
	// /*31*/	'\x00': {I, I, I, I, I, I, I, I}, /* SHIFT_DOWN */
}

var baudotChars = []baudotBits{
	{O, I, I, O, O, O, I, I}, /* A */
	{O, I, O, O, I, I, I, I}, /* B */
	{O, O, I, I, I, O, I, I}, /* C */
	{O, I, O, O, I, O, I, I}, /* D */
	{O, I, O, O, O, O, I, I}, /* E / 3 */
	{O, I, O, I, I, O, I, I}, /* F */
	{O, O, I, O, I, I, I, I}, /* G */
	{O, O, O, I, O, I, I, I}, /* H */
	{O, O, I, I, O, O, I, I}, /* I  / 8 */
	{O, I, I, O, I, O, I, I}, /* J */
	{O, I, I, I, I, O, I, I}, /* K */
	{O, O, I, O, O, I, I, I}, /* L */
	{O, O, O, I, I, I, I, I}, /* M / . */
	{O, O, O, I, I, O, I, I}, /* N */
	{O, O, O, O, I, I, I, I}, /* O / 9 */
	{O, O, I, I, O, I, I, I}, /* P / O */
	{O, I, I, I, O, I, I, I}, /* Q / I  */
	{O, O, I, O, I, O, I, I}, /* R / 4 */
	{O, I, O, I, O, O, I, I}, /* S */
	{O, O, O, O, O, I, I, I}, /* T / 5 */
	{O, I, I, I, O, O, I, I}, /* U / 7 */
	{O, O, I, I, I, I, I, I}, /* V */
	{O, I, I, O, O, I, I, I}, /* W / 2 */
	{O, I, O, I, I, I, I, I}, /* X / / */
	{O, I, O, I, O, I, I, I}, /* Y / 6 */
	{O, I, O, O, O, I, I, I}, /* Z */
	{O, O, O, O, O, O, I, I}, /* NULL */
	{O, O, I, O, O, O, I, I}, /* LF */
	{O, O, O, I, O, O, I, I}, /* SPACE */
	{O, O, O, O, I, O, I, I}, /* CR */
	{O, I, I, O, I, I, I, I}, /* SHIFT_UP */
	{O, I, I, I, I, I, I, I}, /* SHIFT_DOWN */
	{O, O, O, O, O, O, O, O}, /* Open */
	{I, I, I, I, I, I, I, I}, /* closed */
}

var ascii2Punctuation = map[rune]int{
	'-':  CHAR_DASH,
	'?':  CHAR_QUESTION,
	':':  CHAR_COLON,
	'$':  CHAR_DOLLAR,
	'\a': CHAR_BELL,
	'\'': CHAR_APOSTROPHE,
	'`':  CHAR_APOSTROPHE,
	'(':  CHAR_LPHAREN,
	')':  CHAR_RPHAREN,
	'.':  CHAR_PERIOD,
	',':  CHAR_COMMA,
	';':  CHAR_SEMICOLON,
	'/':  CHAR_SOLIDUS,
	'"':  CHAR_QUOTE,
	'\r': CHAR_CR,
	'\n': CHAR_LF,

	'0': CHAR_0,
	'1': CHAR_1,
	'2': CHAR_2,
	'3': CHAR_3,
	'4': CHAR_4,
	'5': CHAR_5,
	'6': CHAR_6,
	'7': CHAR_7,
	'8': CHAR_8,
	'9': CHAR_9,
}

var ascii2Characters = map[rune]int{
	'a': CHAR_A + 0,
	'A': CHAR_A + 0,
	'b': CHAR_A + 1,
	'B': CHAR_A + 1,
	'c': CHAR_A + 2,
	'C': CHAR_A + 2,
	'd': CHAR_A + 3,
	'D': CHAR_A + 3,
	'e': CHAR_A + 4,
	'E': CHAR_A + 4,
	'f': CHAR_A + 5,
	'F': CHAR_A + 5,
	'g': CHAR_A + 6,
	'G': CHAR_A + 6,
	'h': CHAR_A + 7,
	'H': CHAR_A + 7,
	'i': CHAR_A + 8,
	'I': CHAR_A + 8,
	'j': CHAR_A + 9,
	'J': CHAR_A + 9,
	'k': CHAR_A + 10,
	'K': CHAR_A + 10,
	'l': CHAR_A + 11,
	'L': CHAR_A + 11,
	'm': CHAR_A + 12,
	'M': CHAR_A + 12,
	'n': CHAR_A + 13,
	'N': CHAR_A + 13,
	'o': CHAR_A + 14,
	'O': CHAR_A + 14,
	'p': CHAR_A + 15,
	'P': CHAR_A + 15,
	'q': CHAR_A + 16,
	'Q': CHAR_A + 16,
	'r': CHAR_A + 17,
	'R': CHAR_A + 17,
	's': CHAR_A + 18,
	'S': CHAR_A + 18,
	't': CHAR_A + 19,
	'T': CHAR_A + 19,
	'u': CHAR_A + 20,
	'U': CHAR_A + 20,
	'v': CHAR_A + 21,
	'V': CHAR_A + 21,
	'w': CHAR_A + 22,
	'W': CHAR_A + 22,
	'x': CHAR_A + 23,
	'X': CHAR_A + 23,
	'y': CHAR_A + 24,
	'Y': CHAR_A + 24,
	'z': CHAR_A + 25,
	'Z': CHAR_A + 25,
	' ': CHAR_SPACE,
}

// Convert input character to intermediate Baudot representation, a
// numeric value in the rate 0-31.
func intValues(v rune, c *convert) ([]int, bool) {
	var retValues = make([]int, 0, 16)

	// Convert punction character first
	if val, ok := ascii2Punctuation[v]; ok {
		if val == CHAR_LF {
			// newline '\n' is found
			if c.shift {
				c.shift = false
				retValues = append(retValues, CHAR_SHIFT_DOWN)
			}
			retValues = append(retValues, CHAR_LF)
			retValues = append(retValues, CHAR_CR)
		} else {
			if !c.shift {
				c.shift = true
				retValues = append(retValues, CHAR_SHIFT_UP)
			}
			retValues = append(retValues, val)
		}
	} else if val, ok = ascii2Characters[v]; ok {
		if c.shift && val != CHAR_SPACE {
			c.shift = false
			retValues = append(retValues, CHAR_SHIFT_DOWN)
		}
		retValues = append(retValues, val)
	}

	if len(retValues) == 0 {
		return nil, false
	}
	return retValues, true
}

func asciiToBaudotNext(r rune, c *convert) ([]baudotBits, bool) {
	var retBits = make([]baudotBits, 0, 16)

	// Deal with control characters first
	switch r {
	case lineFeed, carriageReturn:
		retBits = append(retBits, baudotConv[carriageReturn])
		retBits = append(retBits, baudotConv[lineFeed])
		return append(retBits, baudotConv[shiftDown]), true
	case spaceCharacter:
		return append(retBits, baudotConv[spaceCharacter]), true
	}

	// Get raw Baudot value
	bits, ok := baudotConv[r]
	if !ok {
		return nil, false
	}

	shift := bits[0]
	if shift != c.shift {
		c.shift = shift
		if shift {
			retBits = append(retBits, baudotConv[shiftUp])
		} else {
			retBits = append(retBits, baudotConv[shiftDown])
		}
		retBits = append(retBits, bits)
	} else {
		return append(retBits, bits), true
	}
	return nil, false
}

func asciiToBaudot(v int) baudotBits {
	if v > len(baudotChars) {
		return baudotChars[CHAR_CLOSED]
	}
	return baudotChars[v]
}

func writeBitsNext(bitsSlice []baudotBits) {
	for _, bits := range bitsSlice {
		for _, bit := range bits[1:] {
			wp.WriteBit(bit)
			wp.DelayMicroseconds(BAUD_DELAY_45)
		}
		wp.DelayMicroseconds(BAUD_DELAY_45 / 2)
	}
}

func writeBits(bits baudotBits) {
	for _, bit := range bits {
		wp.WriteBit(bit)
		wp.DelayMicroseconds(BAUD_DELAY_45)
	}
}

func printRune(r rune, c *convert) {
	if vals, ok := intValues(r, c); ok {
		fmt.Printf("%c", unicode.ToUpper(r))
		for _, val := range vals {
			bits := asciiToBaudot(val)
			writeBits(bits)
		}
	}
}

func printRuneNext(r rune, c *convert) {
	bitsSlice, ok := asciiToBaudotNext(r, c)
	if !ok {
		return
	}
	writeBitsNext(bitsSlice)
}

func printRaw(val int, c *convert) {
	bits := asciiToBaudot(val)
	writeBits(bits)
}

func printRawNext(r rune, c *convert) {
	bits, ok := asciiToBaudotNext(r, c)
	if !ok {
		return
	}
	writeBitsNext(bits)
}

func initializeTeletype(c *convert) {
	// Initialize teletype
	printRaw(CHAR_NULL, c)
	printRaw(CHAR_NULL, c)
	printRaw(CHAR_SHIFT_DOWN, c)
	printRaw(CHAR_CR, c)
	printRaw(CHAR_LF, c)
}

func initializeTeletypeNext(c *convert) {
	printRawNext(nullChar, c)
	printRawNext(nullChar, c)
	printRawNext(shiftDown, c)
	printRawNext(carriageReturn, c)
	printRawNext(lineFeed, c)
}
