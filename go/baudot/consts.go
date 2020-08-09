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

// ASCII to Baudot constants
const (
	CHAR_A          = 0
	CHAR_Z          = 25
	CHAR_NULL       = 26
	CHAR_LF         = 27
	CHAR_SPACE      = 28
	CHAR_CR         = 29
	CHAR_SHIFT_UP   = 30
	CHAR_SHIFT_DOWN = 31
	CHAR_OPEN       = 32
	CHAR_CLOSED     = 33

	CHAR_0 = 15
	CHAR_1 = 16
	CHAR_2 = 22
	CHAR_3 = 4
	CHAR_4 = 17
	CHAR_5 = 19
	CHAR_6 = 24
	CHAR_7 = 20
	CHAR_8 = 8
	CHAR_9 = 14

	CHAR_DASH       = 0
	CHAR_QUESTION   = 1
	CHAR_COLON      = 2
	CHAR_DOLLAR     = 3
	CHAR_BELL       = 6
	CHAR_APOSTROPHE = 9
	CHAR_LPHAREN    = 10
	CHAR_RPHAREN    = 11
	CHAR_PERIOD     = 12
	CHAR_COMMA      = 13
	CHAR_SEMICOLON  = 21
	CHAR_SOLIDUS    = 23
	CHAR_QUOTE      = 25
)
