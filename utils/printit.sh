#!/bin/sh

if [ -z "$1" ]; then
  echo usage: $0 filename
  exit 1
fi

rtty --output-dev message --input-file $1

sox -r 8000 -t raw -u -b message out.wav
