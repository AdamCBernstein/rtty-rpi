#!/bin/bash
export PATH=$PATH:/usr/local/bin
LOG="/tmp/nws-teletype.log"

# Download local area weather forecast, and print on Model 28 Teletype
# This is intented to be run as a cron job.
# Example crontab entry
# */20 * * * * /usr/local/bin/nws-teletype.sh
# This is script is dependent this project
#    https://github.com/AdamCBernstein/nws-http

if [ ! -d "/var/tmp/nws" ]; then
  mkdir "/var/tmp/nws"
fi

if [ -L "/var/tmp/nws.lock" ]; then
  echo "$(date) $0 is already running; 	quitting." >> "$LOG"
  exit 0
fi
ln -s "/var/tmp/nws" "/var/tmp/nws.lock"

file1=/var/tmp/nws/sew.txt
file2=/var/tmp/nws/sew2.txt
file_print=$file1
print=0

if [ ! -f "$file1" ]; then
  file_print="$file1"
  nws-http sew > "$file_print"
  print=1
else
  file_print="$file2"
  nws-http sew > "$file_print"
  cmp "$file1" "$file2"
  if [ $? -ne 0 ]; then
    print=1
    cp "$file2" "$file1"
  fi
fi

if [ $print -eq 1 ]; then
  echo "$(date) rtty $file_print" >> "$LOG"
  nice -n 0 rtty "$file_print" > /dev/null 2>&1
else
  echo "$(date) '$file_print' has not changed" >> "$LOG"
fi
rm -f "/var/tmp/nws.lock"
