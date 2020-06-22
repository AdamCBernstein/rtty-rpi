#include <wiringPi.h>
#define DELAY_STATE 44
int main (void)
{
  wiringPiSetup () ;
  pinMode (0, OUTPUT) ;
  for (;;)
  {
    digitalWrite (0, HIGH) ; delay (DELAY_STATE) ;
    digitalWrite (0,  LOW) ; delay (DELAY_STATE) ;
  }
  return 0 ;
}
