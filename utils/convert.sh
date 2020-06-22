#!/bin/sh
#./rtty --speed 44100  --bits 16 --input-file tst3.txt --output-dev rtty_44100.pcm

#sox -r 8000 -t raw -u -b message out.wav
sox -r 44100 -c 1 -t raw -s -w rtty_44100.pcm rtty_44100.wav
