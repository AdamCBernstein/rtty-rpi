package baudot

import (
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

// BCM Controller -> Raspberry Pi A & B pin mapping
// https://github.com/stianeikeland/go-rpio/blob/master/rpio.go#L35-L60
// https://github.com/stianeikeland/go-rpio.git
//             Rev 2 and 3 Raspberry Pi                        Rev 1 Raspberry Pi (legacy)
//   +-----+---------+----------+---------+-----+      +-----+--------+----------+--------+-----+
//   | BCM |   Name  | Physical | Name    | BCM |      | BCM | Name   | Physical | Name   | BCM |
//   +-----+---------+----++----+---------+-----+      +-----+--------+----++----+--------+-----+
//   |     |    3.3v |  1 || 2  | 5v      |     |      |     | 3.3v   |  1 ||  2 | 5v     |     |
//   |   2 |   SDA 1 |  3 || 4  | 5v      |     |      |   0 | SDA    |  3 ||  4 | 5v     |     |
//   |   3 |   SCL 1 |  5 || 6  | 0v      |     |      |   1 | SCL    |  5 ||  6 | 0v     |     |
//   |   4 | GPIO  7 |  7 || 8  | TxD     | 14  |      |   4 | GPIO 7 |  7 ||  8 | TxD    |  14 |
//   |     |      0v |  9 || 10 | RxD     | 15  |      |     | 0v     |  9 || 10 | RxD    |  15 |
//   |  17 | GPIO  0 | 11 || 12 | GPIO  1 | 18  |      |  17 | GPIO 0 | 11 || 12 | GPIO 1 |  18 |
//   |  27 | GPIO  2 | 13 || 14 | 0v      |     |      |  21 | GPIO 2 | 13 || 14 | 0v     |     |
//   |  22 | GPIO  3 | 15 || 16 | GPIO  4 | 23  |      |  22 | GPIO 3 | 15 || 16 | GPIO 4 |  23 |
//   |     |    3.3v | 17 || 18 | GPIO  5 | 24  |      |     | 3.3v   | 17 || 18 | GPIO 5 |  24 |
//   |  10 |    MOSI | 19 || 20 | 0v      |     |      |  10 | MOSI   | 19 || 20 | 0v     |     |
//   |   9 |    MISO | 21 || 22 | GPIO  6 | 25  |      |   9 | MISO   | 21 || 22 | GPIO 6 |  25 |
//   |  11 |    SCLK | 23 || 24 | CE0     | 8   |      |  11 | SCLK   | 23 || 24 | CE0    |   8 |
//   |     |      0v | 25 || 26 | CE1     | 7   |      |     | 0v     | 25 || 26 | CE1    |   7 |
//   |   0 |   SDA 0 | 27 || 28 | SCL 0   | 1   |      +-----+--------+----++----+--------+-----+
//   |   5 | GPIO 21 | 29 || 30 | 0v      |     |
//   |   6 | GPIO 22 | 31 || 32 | GPIO 26 | 12  |
//   |  13 | GPIO 23 | 33 || 34 | 0v      |     |
//   |  19 | GPIO 24 | 35 || 36 | GPIO 27 | 16  |
//   |  26 | GPIO 25 | 37 || 38 | GPIO 28 | 20  |
//   |     |      0v | 39 || 40 | GPIO 29 | 21  |
//   +-----+---------+----++----+---------+-----+
var (
	GPIO0 = 17
	GPIO1 = 18
	GPIO2 = 27
	GPIO3 = 22
	GPIO4 = 23
	GPIO5 = 24
	GPIO6 = 25
	GPIO7 = 4
)

type Convert struct {
	bitTimeDuration time.Duration
	shift           bool
	pin             rpio.Pin
	motor           rpio.Pin
}

func New(speed int) (*Convert, error) {
	if err := rpio.Open(); err != nil {
		return nil, err
	}
	c := &Convert{bitTimeDuration: time.Duration(speed),
		shift: false,
		pin:   rpio.Pin(GPIO0),
		motor: rpio.Pin(GPIO1)}
	c.pin.Output()
	c.motor.Output()
	c.initializeTeletype()
	return c, nil
}

func (c *Convert) MotorControl(enable bool) {
	if enable {
		c.motor.High()
	} else {
		c.motor.Low()
	}
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
