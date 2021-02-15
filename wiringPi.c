#ifndef _WIRINGPI_H_
#define _WIRINGPI_H_ 1

#include <errno.h>
#include <stdio.h>

#define HIGH 1
#define LOW  0
#define OUTPUT 0

static FILE *GPIO_pins[8]; /* Support only GPIO0-GPIO-7 */

int wiringPiSetup(void)
{
    errno = ENOTSUP;
    return -1;
}

void pinMode(int pin, int mode)
{
    char name[64];
    if (mode == OUTPUT && !GPIO_pins[pin])
    {
        sprintf(name, "/tmp/GPIO_%d", pin);
        GPIO_pins[pin] = fopen(name, "w+");
    }
}

void digitalWrite(int pin, int onOff)
{
    static int count[8];
    static int groups[8];

    if (GPIO_pins[pin])
    {
        fprintf(GPIO_pins[pin], "%d", onOff);
        count[pin]++;

        if (count[pin] == 8)
        {
            if (groups[pin] == 8)
            {
                fprintf(GPIO_pins[pin], "\n");
                groups[pin] = 0;
            }
            else
            {
                fprintf(GPIO_pins[pin], " ");
                groups[pin]++;
            }
            count[pin] = 0;
        }
        fflush(GPIO_pins[pin]);
    }
}

void  delay(int usec)
{
}


#endif
