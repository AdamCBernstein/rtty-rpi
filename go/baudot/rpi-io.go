package baudot

import (
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

var (
	// Use mcu pin 17, which corresponds to logical GPIO0 pin
	pinGpio = rpio.Pin(17)
)

type Convert struct {
	bitTimeDuration time.Duration
	shift           bool
	pin             rpio.Pin
}

func New() *Convert {
	if err := rpio.Open(); err != nil {
		return nil
	}
	pinGpio.Output()
	c := &Convert{bitTimeDuration: BAUD_DELAY_45, shift: false, pin: pinGpio}
	c.initializeTeletype()
	return c
}

func (c *Convert) WriteBits(bitsChar []baudotBits) {
	for _, bits := range bitsChar {
		for _, bit := range bits[:LTRS_FIGS_BIT] {
			if !bit {
				c.pin.High()
			} else {
				c.pin.Low()
			}
			time.Sleep(time.Microsecond * c.bitTimeDuration)
		}
		time.Sleep((time.Microsecond * c.bitTimeDuration * 50) / 100)
	}
}
