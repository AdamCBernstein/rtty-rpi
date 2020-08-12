package main

import (
	"bufio"
	"os"

	"rtty/baudot"
)

func printTest(baudot *baudot.Convert) {
	spaces := "SpaceTest:                                                        "
	I0 := "Z0Y1X2W3V4U5T6S7R8P9Q0Z0Y1X2W3V4U5T6S7R8P9Q0Z0Y1X2W3V4U5T6S7R8P9Q0"
	test := "the quick brown fox jumped over the lazy dog's back     1234567890\n" +
		"ryryryryryryryryryryryryryryryryryryryryrryyryryryryryryryryryryry\n" +
		"sgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsgsg"

	baudot.Print(spaces)
	baudot.Print(I0)
	baudot.Print("\n")
	baudot.Print(test)
	baudot.Print(test)
	baudot.Print(test)
	baudot.Print(test)
	baudot.Print("\n\n")
}

func printFile(fname string, b *baudot.Convert) error {
	fp, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fp.Close()

	scan := bufio.NewScanner(fp)
	for scan.Scan() {
		b.Print(scan.Text())
	}
	return nil
}

func main() {
	// Print named file, if one is provided. Otherwise, print test
	b := baudot.New()
	if len(os.Args) > 1 {
		printFile(os.Args[1], b)
	} else {
		printTest(b)
	}
}
