#!/bin/bash
processed_log="/var/tmp/usgs/printed_files.log"

function trap_error {
  local sts="$?"
  if [ $sts -ne 0 ]; then
    echo "FATAL: error=$sts"
  fi
  exit 1
}

function init() {
  set -eou pipefail
  touch "$processed_log"
}

function print_teletype() {
  local file="$1"
  rtty "$file"
}

function extract_text() {
  local file="$1"
  local report_file="usgs-report.txt"
  local report_latest="usgs-report-latest.txt"
  local report_subject=$(cat "$file" | grep "^Subject: ")
  local subject=""
  local subject_title=""

  cat "$file"| \
    sed -n '\|Content-Type: text/plain|,\|Content-Type: text/html|p' | \
    tail -n +4 | \
    sed  '\|^This is a computer-generated message and has not |,$d' > "$report_file"
    rm -f "$file"

  # Subject: 2016-09-01 16:37:58 UPDATED: (M7.1) OFF EAST COAST OF THE NORTH ISLAND, N.Z. -37.5 179.2 (39068)

  subject=$(echo "$report_subject" | awk '{print $2 " " $3}')
  subject="(USGS) ENG - Preliminary Earthquake Report ${subject}"
  subject_title=$(echo $report_subject | sed 's/.............................//')
  if [ ! -f "$report_latest" ]; then
    echo "$subject"        > "$report_latest"
    echo "$subject_title" >> "$report_latest"
    cat "$report_file"    >> "$report_latest"
    print_teletype "$report_latest"
    echo "$file" >> "$processed_log"
  else
    # Verify current file has not been seen before
    if [ $(grep -c "$file" "$processed_log") -eq 0 ]; then
      echo "$file" >> "$processed_log"
      set +e
      cmp -s "$report_file" "$report_latest"
      if [ $? -ne 0 ]; then
        echo "$subject"        > "$report_latest"
        echo "$subject_title" >> "$report_latest"
        cat "$report_file"    >> "$report_latest"
        print_teletype "$report_latest"
      fi
      set -e
    fi
  fi
}

function print_files() {
  local files=$(find . -maxdepth 1 -name '*.robin:*' | sort -n)
  for f in $files; do
    extract_text "$f"
  done
}

function main() {
  init

  trap trap_error 1 2 3 15

  print_files
}

main
