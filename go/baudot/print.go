package baudot

import (
	"rtty/gpio"
)

func New() *convert {
	wiringpi.Initialize(wiringpi.GPIO_PIN0)
	c := &convert{shift: false}
	initializeTeletype(c)
	return c
}

func (c *convert) Print(line string) {
	column := 0
	for _, char := range line {
		printRune(char, c)
		column++

		if column > COLUMN_MAX {
			printRune('\n', c)
			column = 0
		}
		if char == '\n' {
			column = 0
		}
	}
	printRune('\n', c)
}
