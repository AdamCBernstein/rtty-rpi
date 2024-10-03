package baudot

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func (c *Convert) initializeTeletype() {
	c.printRune(shiftDown)
	c.printRune(shiftDown)
	c.printRune(carriageReturn)
	c.printRune(lineFeed)
	c.printRune(shiftDown)
}

func (c *Convert) Write(line []byte) (n int, err error) {
	c.PrintLine(string(line))
	return len(line), nil
}

func (c *Convert) PrintLine(line string) {
	column := 0
	for _, char := range line {
		c.printRune(char)
		column++

		if column > COLUMN_MAX {
			c.printRune('\n')
			column = 0
		}
		if char == '\n' {
			column = 0
		}
	}
	c.printRune('\n')
}

func (c *Convert) PrintRune(r rune) {
	c.printRune(r)
}

func (c *Convert) printRune(r rune) {
	rUpper := unicode.ToUpper(r)
	bitsSlice, ok := c.asciiToBaudot(rUpper)
	if !ok {
		return
	}

	// Give some console feedback. TODO: Make callback function
	fmt.Printf("%c", rUpper)
	c.WriteBits(bitsSlice)
}

func (c *Convert) PrintFile(fname string) error {
	if fp, err := os.Open(fname); err != nil {
		return err
	} else {
		defer fp.Close()

		scan := bufio.NewScanner(fp)
		for scan.Scan() {
			c.PrintLine(scan.Text())
		}
	}
	return nil
}
