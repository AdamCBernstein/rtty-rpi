#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <wiringPi.h>
#define DELAY_STATE 44
int main (int argc, char *argv[])
{
  int pin = 0;
  int wait = DELAY_STATE;

  if (argc > 1 && (!strcmp(argv[1], "-h") || !strcmp(argv[1], "--help"))) {
    printf("usage: blink [GPIO_PIN] [Delay ms] [0 (off) | 1 (on)]\n");
    printf("             These options stack; later options require previous values\n");
    return 0;
  }
  if (argc > 1)
  {
    pin = atoi(argv[1]);
  }
  if (argc > 2)
  {
    wait = atoi(argv[2]);
  }

  wiringPiSetup () ;
  pinMode (pin, OUTPUT) ;

  if (argc > 3)
  {
    /* Set "pin" value to LOW (off) or HIGH (on) */
    if (atoi(argv[3]) == 0) 
    {
      digitalWrite (pin, LOW);
    }
    else
    {
      digitalWrite (pin, HIGH);
    }
    return 0;
  }

  for (;;)
  {
    digitalWrite (pin, HIGH); delay (wait);
    digitalWrite (pin,  LOW); delay (wait);
  }
  return 0 ;
}
