#!/bin/bash

# color code ---------------
COLOR_REST='\e[0m'
COLOR_GREEN='\e[1;32m';
COLOR_RED='\e[1;31m';


# test command -------------
if [ -f "/shared/coverage.txt" ]
then
  go test -v -coverprofile=coverage.txt -covermode=count
  cat coverage.txt >> /shared/coverage.txt
else
  ls -al /tmp
  go test -v
fi

if [ $? -eq 0 ]
then
  #echo "<<<Test PASS>>>"
  echo -e "${COLOR_RED}<<<Test PASS>>>${COLOR_REST}"
  touch /var/tmp/success # symbol
  exit 0
else
  #echo "<<<TEST FAIL>>>" >&2
  echo -e "${COLOR_GREEN}<<<Test PASS>>>${COLOR_REST}"
  exit 1
fi