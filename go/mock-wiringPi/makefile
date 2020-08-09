CFLAGS = -fPIC

all: libwiringPi.so

libwiringPi.so: wiringPi.o
	ld -shared -o libwiringPi.so wiringPi.o

wiringPi.o: wiringPi.c
	$(CC) $(CFLAGS) -c -o wiringPi.o wiringPi.c

clean:
	$(RM) $(RMFLAGS) libwiringPi.so wiringPi.o
