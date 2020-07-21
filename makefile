CFLAGS = -I/usr/local/include -Wall -g
all: rtty-pi

rtty-pi: hack-pi.o 
	cc -o rtty-pi hack-pi.o -L/usr/local/lib -lwiringPi
	sudo chown root rtty-pi
	sudo chmod 4755 rtty-pi
clean: 
	rm -f hack-pi.o rtty-pi
