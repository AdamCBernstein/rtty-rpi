package baudot

// Baud speed durations
const (
	BAUD_DELAY_45 = 22000 /* 22ms = 45 baud, 60WPM; in uSeconds */
	BAUD_DELAY_50 = 20000 /* 20ms = 50 baud, 66WPM; in uSeconds*/
	BAUD_DELAY_74 = 13470 /* 13ms = 74 baud, 100WPM; in uSeconds*/
)

// Misc constants
const (
	COLUMN_MAX = 76
)

const I = true
const O = false
const shiftUp = '\x02'
const shiftDown = '\x01'
const lineFeed = '\n'
const carriageReturn = '\r'
const spaceCharacter = ' '
const nullChar = '\xfe'
const LETTERS = 0x01
const FIGURES = 0x02
const LTRS_FIGS_BIT = 7

// ITA-2 / US-TTY table. Reference:
// https://en.wikipedia.org/wiki/Baudot_code#ITA_2_and_US-TTY
//
// format: bit0-bit6 = | start | data-5-bits | stop |
//         bit7 = up/down shift
// S = START-BIT; E=END-BIT; L=LETTERS F=FIGURES
// | S | D1 | D2 | D3 | D4 | D5 |  E  | L/F |
var baudotConv = map[rune]baudotBits{
	// All entries here are LETTERS; bit7 = O
	'A':     {O, I, I, O, O, O, I, O}, /* A / - */
	'B':     {O, I, O, O, I, I, I, O}, /* B / ? */
	'C':     {O, O, I, I, I, O, I, O}, /* C / : $ */
	'D':     {O, I, O, O, I, O, I, O}, /* D */
	'E':     {O, I, O, O, O, O, I, O}, /* E / 3 */
	'F':     {O, I, O, I, I, O, I, O}, /* F */
	'G':     {O, O, I, O, I, I, I, O}, /* G / \a */
	'H':     {O, O, O, I, O, I, I, O}, /* H / 8 */
	'I':     {O, O, I, I, O, O, I, O}, /* I  / 8 */
	'J':     {O, I, I, O, I, O, I, O}, /* J / ' */
	'K':     {O, I, I, I, I, O, I, O}, /* K / ( */
	'L':     {O, O, I, O, O, I, I, O}, /* L / ) */
	'M':     {O, O, O, I, I, I, I, O}, /* M / . */
	'N':     {O, O, O, I, I, O, I, O}, /* N / , */
	'O':     {O, O, O, O, I, I, I, O}, /* 0  / 9 */
	'P':     {O, O, I, I, O, I, I, O}, /* P / O */
	'Q':     {O, I, I, I, O, I, I, O}, /* Q / 1 */
	'R':     {O, O, I, O, I, O, I, O}, /* R / 4 */
	'S':     {O, I, O, I, O, O, I, O}, /* S */
	'T':     {O, O, O, O, O, I, I, O}, /* T / 5 */
	'U':     {O, I, I, I, O, O, I, O}, /* U / 7 */
	'V':     {O, O, I, I, I, I, I, O}, /* V / ; */
	'W':     {O, I, I, O, O, I, I, O}, /* W / 2 */
	'X':     {O, I, O, I, I, I, I, O}, /* X / / */
	'Y':     {O, I, O, I, O, I, I, O}, /* Y / 6 */
	'Z':     {O, I, O, O, O, I, I, O}, /* Z / " */
	'\xfe':  {O, O, O, O, O, O, I, O}, /* NULL */
	'\n':    {O, O, I, O, O, O, I, O}, /* LF */
	' ':     {O, O, O, I, O, O, I, O}, /* SPACE */
	'\r':    {O, O, O, O, I, O, I, O}, /* CR */
	FIGURES: {O, I, I, O, I, I, I, O}, /* SHIFT_UP */
	LETTERS: {O, I, I, I, I, I, I, O}, /* SHIFT_DOWN */

	// All entries here are FIGURES; bit7 = 1
	'-': {O, I, I, O, O, O, I, I}, /* A / - */
	'?': {O, I, O, O, I, I, I, I}, /* B / ? */
	':': {O, O, I, I, I, O, I, I}, /* C / : */
	'$': {O, I, O, O, I, O, I, I}, /* D / $ */
	'3': {O, I, O, O, O, O, I, I}, /* E / 3 */
	//        BEL 'A': {O, I, O, I, I, O, I, I}, /* F */
	'\a': {O, O, I, O, I, I, I, I}, /* G / \a */
	//        UNKNOWN {O, O, O, I, O, I, I, I}, /* H */
	'8':  {O, O, I, I, O, O, I, I}, /* I  / 8 */
	'\'': {O, I, I, O, I, O, I, I}, /* J / ' */
	'(':  {O, I, I, I, I, O, I, I}, /* K / ( */
	')':  {O, O, I, O, O, I, I, I}, /* L / ) */
	'.':  {O, O, O, I, I, I, I, I}, /* M / . */
	',':  {O, O, O, I, I, O, I, I}, /* N / , */
	'9':  {O, O, O, O, I, I, I, I}, /* 0  / 9 */
	'0':  {O, O, I, I, O, I, I, I}, /* P / O */
	'1':  {O, I, I, I, O, I, I, I}, /* Q / 1 */
	'4':  {O, O, I, O, I, O, I, I}, /* R / 4 */
	// ALT-BEL  {O, I, O, I, O, O, I, I}, /* S */
	'5': {O, O, O, O, O, I, I, I}, /* T / 5 */
	'7': {O, I, I, I, O, O, I, I}, /* U / 7 */
	';': {O, O, I, I, I, I, I, I}, /* V / ( */
	'2': {O, I, I, O, O, I, I, I}, /* W / 2 */
	'/': {O, I, O, I, I, I, I, I}, /* X / / */
	'6': {O, I, O, I, O, I, I, I}, /* Y / 6 */
	'"': {O, I, O, O, O, I, I, I}, /* Z / " */
	// What does upper case control characters even mean?
	//      nullChar: {O, O, O, O, O, O, I, I}, /* NULL */
	//      lineFeed: {O, O, I, O, O, O, I, I}, /* LF */
	//      spaceChar: {O, O, O, I, O, O, I, I}, /* SPACE */
	//      carriageReturn: {O, O, O, O, I, O, I, I}, /* CR */
	//      {O, I, I, O, I, I, I, I}, /* SHIFT_UP */
	//      {O, I, I, I, I, I, I, I}, /* SHIFT_DOWN */
}
