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

// Define _true so it is as long as "false" to keep below table aligned
const _true = true

var baudotChars = []baudotBits{
	{false, _true, _true, false, false, false, _true, _true}, /* A */
	{false, _true, false, false, _true, _true, _true, _true}, /* B */
	{false, false, _true, _true, _true, false, _true, _true}, /* C */
	{false, _true, false, false, _true, false, _true, _true}, /* D */
	{false, _true, false, false, false, false, _true, _true}, /* E / 3 */
	{false, _true, false, _true, _true, false, _true, _true}, /* F */
	{false, false, _true, false, _true, _true, _true, _true}, /* G */
	{false, false, false, _true, false, _true, _true, _true}, /* H */
	{false, false, _true, _true, false, false, _true, _true}, /* I  / 8 */
	{false, _true, _true, false, _true, false, _true, _true}, /* J */
	{false, _true, _true, _true, _true, false, _true, _true}, /* K */
	{false, false, _true, false, false, _true, _true, _true}, /* L */
	{false, false, false, _true, _true, _true, _true, _true}, /* M / . */
	{false, false, false, _true, _true, false, _true, _true}, /* N */
	{false, false, false, false, _true, _true, _true, _true}, /* O / 9 */
	{false, false, _true, _true, false, _true, _true, _true}, /* P / false */
	{false, _true, _true, _true, false, _true, _true, _true}, /* Q / _true  */
	{false, false, _true, false, _true, false, _true, _true}, /* R / 4 */
	{false, _true, false, _true, false, false, _true, _true}, /* S */
	{false, false, false, false, false, _true, _true, _true}, /* T / 5 */
	{false, _true, _true, _true, false, false, _true, _true}, /* U / 7 */
	{false, false, _true, _true, _true, _true, _true, _true}, /* V */
	{false, _true, _true, false, false, _true, _true, _true}, /* W / 2 */
	{false, _true, false, _true, _true, _true, _true, _true}, /* X / / */
	{false, _true, false, _true, false, _true, _true, _true}, /* Y / 6 */
	{false, _true, false, false, false, _true, _true, _true}, /* Z */
	{false, false, false, false, false, false, _true, _true}, /* NULL */
	{false, false, _true, false, false, false, _true, _true}, /* LF */
	{false, false, false, _true, false, false, _true, _true}, /* SPACE */
	{false, false, false, false, _true, false, _true, _true}, /* CR */
	{false, _true, _true, false, _true, _true, _true, _true}, /* SHIFT_UP */
	{false, _true, _true, _true, _true, _true, _true, _true}, /* SHIFT_DOWN */
	{false, false, false, false, false, false, false, false}, /* Open */
	{_true, _true, _true, _true, _true, _true, _true, _true}, /* closed */
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

func asciiToBaudot(v int) baudotBits {
	if v > len(baudotChars) {
		return baudotChars[CHAR_CLOSED]
	}
	return baudotChars[v]
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

func printRaw(val int, c *convert) {
	bits := asciiToBaudot(val)
	writeBits(bits)
}

func initializeTeletype(c *convert) {
	// Initialize teletype
	printRaw(CHAR_NULL, c)
	printRaw(CHAR_NULL, c)
	printRaw(CHAR_SHIFT_DOWN, c)
	printRaw(CHAR_CR, c)
	printRaw(CHAR_LF, c)
}
