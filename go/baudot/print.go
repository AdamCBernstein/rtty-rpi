package baudot

import (
	wp "rtty/gpio"
	"time"
)

func New() *convert {
	wp.Initialize(wp.GPIO_PIN0)
	time.Sleep(time.Second)

	c := &convert{
		shift: false,
	}
	initializeTeletype(c)

	return c
}

func (c *convert) Print(line string) {
	i := 0
	for _, char := range line {
		printRune(char, c)
		i++

		if i > COLUMN_MAX {
			printRune('\n', c)
			i = 0
		}
		if char == '\n' {
			i = 0
		}
	}
	printRune('\n', c)
}
