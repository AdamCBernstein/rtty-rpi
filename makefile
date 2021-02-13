CFLAGS = -I/usr/local/include -Wall -g
all: rtty-pi go/rtty blink

rtty-pi: hack-pi.o 
	cc -o rtty-pi hack-pi.o -L/usr/local/lib -lwiringPi
	sudo chown root rtty-pi
	sudo chmod 4755 rtty-pi
blink: blink.o 
	cc -o blink blink.o -L/usr/local/lib -lwiringPi
	sudo chown root blink
	sudo chmod 4755 blink

go/rtty: go/baudot/print.go go/baudot/test.go go/baudot/encode.go go/baudot/consts.go go/baudot/rpi-io.go go/main.go
	cd go && go fmt ./... &&  go build &&  sudo chown root rtty && sudo chmod 4755 rtty

clean: 
	rm -f hack-pi.o rtty-pi blink.o blink go/rtty
