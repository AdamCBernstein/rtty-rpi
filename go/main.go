package main

import (
//	"fmt"
	"rtty/baudot"
//	wp "rtty/gpio"
//	"time"
)

func main() {
	test := "the quick brown fox jumped over the lazy dog's back 1234567890\nryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryry"
//	wp.Initialize(wp.GPIO_PIN0)
	baudot := baudot.New()
	baudot.Print(test)
/*
	for _, c := range test {
		vals, ok := baudot.ToIntValues(rune(c))
		if ok {
			fmt.Println(vals)
		}
	}
*/

/*
	for {
		wp.WriteBitOn()
		time.Sleep(time.Millisecond * 100)
		wp.WriteBitOff()
		time.Sleep(time.Millisecond * 100)
	}
*/
}
