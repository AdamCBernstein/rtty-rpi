package wiringpi

// #cgo LDFLAGS: -Wl,-rpath=/usr/lib -L/usr/lib -lwiringPi
// #include <wiringPi.h>
import "C"

func Initialize(pin int) {
	C.wiringPiSetup()
	C.pinMode(C.int(pin), C.OUTPUT)
}

func WriteBitOn() {
	C.digitalWrite(0, C.HIGH)
}

func WriteBitOff() {
	C.digitalWrite(0, C.LOW)
}
