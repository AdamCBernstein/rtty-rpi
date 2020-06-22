CFLAGS = -I/usr/local/include -Wall -g
all: rtty-pi

rtty-pi: hack-pi.o 
	cc -o rtty-pi hack-pi.o -L/usr/local/lib -lwiringPi
clean: 
	rm -f hack-pi.o rtty-pi
