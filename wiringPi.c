#ifndef _WIRINGPI_H_
#define _WIRINGPI_H_ 1

#include <errno.h>

#define HIGH 1
#define LOW  0
#define OUTPUT 0

int wiringPiSetup(void)
{
    errno = ENOTSUP;
    return -1;
}

void pinMode(int pin, int mode)
{
}

void digitalWrite(int pin, int onOff)
{
}

void  delay(int usec)
{
}


#endif
