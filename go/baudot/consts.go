package baudot

// Baud speed durations, in microseconds (uSeconds)
const (
	BAUD_DELAY_45 = 22000 // 22ms = 45 baud, 60WPM; in uSeconds
	BAUD_DELAY_50 = 20000 // 20ms = 50 baud, 66WPM; in uSeconds
	BAUD_DELAY_57 = 18000 // 18ms = 56.9 baud, 75WPM; in uSeconds
	BAUD_DELAY_74 = 13470 // 13ms = 74 baud, 100WPM; in uSeconds
)

// Misc constants
const (
	COLUMN_MAX = 76
)

// Symbolic names for rune (ASCII) values used as keys in Baudot table
const (
	shiftDown      = '\x01'
	shiftUp        = '\x02'
	bellChar       = '\a'
	lineFeed       = '\n'
	carriageReturn = '\r'
	tabChar        = '\t'
	spaceChar      = ' '
	nullChar       = '\xfe'
)

// Constants related to Teletype code conversion
const (
	LETTERS       = 0x01
	FIGURES       = 0x02
	LTRS_FIGS_BIT = 7
)

// Define I ("one"; true), O ("zero"; false) to populate Baudot
// code table below. Makes reading the bool values (bits) easier.
// These are the letters 'I' and 'O', respectively
const I = true
const O = false

// ITA-2 / US-TTY table. Reference:
// https://en.wikipedia.org/wiki/Baudot_code#ITA_2_and_US-TTY
//
// format: bit0-bit6 = | start | D1-D5  | stop |
//         bit7 = LETTERS / FIGURES shift "bit" value
//
// Baudot table format below:
// S = START-BIT; E=END-BIT; D1-D5=Data Bits L=LETTERS F=FIGURES
var baudotConv = map[rune]baudotBits{
	// All entries here are LETTERS; bit7 = O
	//             | S |D1|D2|D3|D4|D5|E| L/F|
	'A':            {O, I, I, O, O, O, I, O}, // A / -
	'B':            {O, I, O, O, I, I, I, O}, // B / ?
	'C':            {O, O, I, I, I, O, I, O}, // C / :
	'D':            {O, I, O, O, I, O, I, O}, // D / $
	'E':            {O, I, O, O, O, O, I, O}, // E / 3
	'F':            {O, I, O, I, I, O, I, O}, // F / !
	'G':            {O, O, I, O, I, I, I, O}, // G / &
	'H':            {O, O, O, I, O, I, I, O}, // H / 8
	'I':            {O, O, I, I, O, O, I, O}, // I  / 8
	'J':            {O, I, I, O, I, O, I, O}, // J / ' (or BELL)
	'K':            {O, I, I, I, I, O, I, O}, // K / (
	'L':            {O, O, I, O, O, I, I, O}, // L / )
	'M':            {O, O, O, I, I, I, I, O}, // M / .
	'N':            {O, O, O, I, I, O, I, O}, // N / ,
	'O':            {O, O, O, O, I, I, I, O}, // 0  / 9
	'P':            {O, O, I, I, O, I, I, O}, // P / O
	'Q':            {O, I, I, I, O, I, I, O}, // Q / 1
	'R':            {O, O, I, O, I, O, I, O}, // R / 4
	'S':            {O, I, O, I, O, O, I, O}, // S / BELL (or ')
	'T':            {O, O, O, O, O, I, I, O}, // T / 5
	'U':            {O, I, I, I, O, O, I, O}, // U / 7
	'V':            {O, O, I, I, I, I, I, O}, // V / ;
	'W':            {O, I, I, O, O, I, I, O}, // W / 2
	'X':            {O, I, O, I, I, I, I, O}, // X / /
	'Y':            {O, I, O, I, O, I, I, O}, // Y / 6
	'Z':            {O, I, O, O, O, I, I, O}, // Z / "
	nullChar:       {O, O, O, O, O, O, I, O}, // NULL
	lineFeed:       {O, O, I, O, O, O, I, O}, // LF
	spaceChar:      {O, O, O, I, O, O, I, O}, // SPACE
	carriageReturn: {O, O, O, O, I, O, I, O}, // CR
	FIGURES:        {O, I, I, O, I, I, I, O}, // SHIFT_UP
	LETTERS:        {O, I, I, I, I, I, I, O}, // SHIFT_DOWN

	// All entries here are FIGURES; bit7 = 1
	'-':      {O, I, I, O, O, O, I, I}, // A / -
	'?':      {O, I, O, O, I, I, I, I}, // B / ?
	':':      {O, O, I, I, I, O, I, I}, // C / :
	'$':      {O, I, O, O, I, O, I, I}, // D / $
	'3':      {O, I, O, O, O, O, I, I}, // E / 3
	'!':      {O, I, O, I, I, O, I, I}, // F / !
	bellChar: {O, O, I, O, I, I, I, I}, // G / BELL
	'8':      {O, O, I, I, O, O, I, I}, // I  / 8
	'\'':     {O, I, I, O, I, O, I, I}, // J / '
	'(':      {O, I, I, I, I, O, I, I}, // K / (
	')':      {O, O, I, O, O, I, I, I}, // L / )
	'.':      {O, O, O, I, I, I, I, I}, // M / .
	',':      {O, O, O, I, I, O, I, I}, // N / ,
	'9':      {O, O, O, O, I, I, I, I}, // 0  / 9
	'0':      {O, O, I, I, O, I, I, I}, // P / O
	'1':      {O, I, I, I, O, I, I, I}, // Q / 1
	'4':      {O, O, I, O, I, O, I, I}, // R / 4
	tabChar:  {O, I, O, I, O, O, I, I}, // S
	'5':      {O, O, O, O, O, I, I, I}, // T / 5
	'7':      {O, I, I, I, O, O, I, I}, // U / 7
	';':      {O, O, I, I, I, I, I, I}, // V / (
	'2':      {O, I, I, O, O, I, I, I}, // W / 2
	'/':      {O, I, O, I, I, I, I, I}, // X / /
	'6':      {O, I, O, I, O, I, I, I}, // Y / 6
	'"':      {O, I, O, O, O, I, I, I}, // Z / "
	// What does upper case control characters even mean?
	//      nullChar: {O, O, O, O, O, O, I, I}, // NULL
	//      lineFeed: {O, O, I, O, O, O, I, I}, // LF
	//      spaceChar: {O, O, O, I, O, O, I, I}, // SPACE
	//      carriageReturn: {O, O, O, O, I, O, I, I}, // CR
	//      {O, I, I, O, I, I, I, I}, // SHIFT_UP
	//      {O, I, I, I, I, I, I, I}, // SHIFT_DOWN
}
