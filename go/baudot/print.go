package baudot

import (
	"rtty/gpio"
)

func New() *Convert {
	wiringpi.Initialize(wiringpi.GPIO_PIN0)
	c := &Convert{shift: false}
	initializeTeletype(c)
	return c
}

func (c *Convert) Print(line string) {
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
