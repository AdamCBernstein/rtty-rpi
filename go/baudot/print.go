package baudot

import (
	"bufio"
	"fmt"
	"os"
)

func (c *Convert) initializeTeletype() {
	fmt.Fprintf(c, "%c", shiftDown)
	fmt.Fprintf(c, "%c", shiftDown)
	fmt.Fprintf(c, "%c", carriageReturn)
	fmt.Fprintf(c, "%c", lineFeed)
	fmt.Fprintf(c, "%c", shiftDown)
}

// Write makes this "interface compatible" with io.Write() and now fmt.Print/fmt.Fprint can be used
func (c *Convert) Write(line []byte) (n int, err error) {
	c.printLine(string(line))
	return len(line), nil
}

// printLine does do more work than Fprintf("%s"), by handling column positioning and sending
// a \r\n character to the Teletype.
func (c *Convert) printLine(line string) {
	column := 0
	for _, char := range line {
		fmt.Fprintf(c, "%c", char)
		column++

		if column > ColumnMax {
			fmt.Fprintf(c, "%c", '\n')
			column = 0
		}
		if char == '\n' {
			column = 0
		}
	}
	fmt.Fprintf(c, "%c", '\n')
}

// PrintFile sends the entire contents of "file" to the teletype
func (c *Convert) PrintFile(fname string) error {
	fp, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fp.Close()

	scan := bufio.NewScanner(fp)
	for scan.Scan() {
		fmt.Fprintf(c, "%s", scan.Text())
	}
	return nil
}
