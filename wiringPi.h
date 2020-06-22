#ifndef _WIRINGPI_H_
#define _WIRINGPI_H_ 1

#define HIGH 1
#define LOW  0
#define OUTPUT 0

int wiringPiSetup(void);
int pinMode(int, int);
int digitalWrite(int, int);
int delay(int);


#endif
