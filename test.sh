#!/bin/bash

## Varaibles
URL=http://localhost:8080/api/mb/tcp
SLAVE=$(docker inspect --format '{{ .NetworkSettings.IPAddress }}' mbweb_slave_1)

## Unit-Testable Shell Scripts (http://eradman.com/posts/ut-shell-scripts.html)
typeset -i tests_run=0
typeset -i fail_run=0
function try { 
  this="$1"
  echo "####" $this
}
trap 'printf "$0: exit code $? on line $LINENO\nFAIL: $this\n"; exit 1' ERR
function assert {
    let tests_run+=1
    [ "$1" = "$2" ] && { echo -n "."; return; }
    printf "\nFAIL: $this\n'$1' != '$2'\n"; exit 1
}
function check_200_status {
    let tests_run+=1
    echo "resp:" $out
    [[ "$1" == 200* ]] && { echo "---------------------------------"; return; }
    printf "@@@FAIL: '$1'\n"; let fail_run+=1; echo "---------------------------------";
}
function check_ok_status {
    let tests_run+=1
    echo "resp:" $out
    [[ "$1" == *'"status":"ok"'* ]] && { echo "---------------------------------"; return; }
    printf "@@@FAIL: '$1'\n"; let fail_run+=1; echo "---------------------------------";
}
function check_not_ok_status {
    let tests_run+=1
    echo "resp:" $out
    [[ "$1" != *'"status":"ok"'* ]] && { echo "---------------------------------"; return; }
    printf "@@@FAIL: '$1'\n"; let fail_run+=1; echo "---------------------------------";
}

# HTTP GET
function GET {
    echo "req :" $1
    out=$(curl -sL $1)
}
## end

###############################################################
echo "---------------------------------" # start

try "Test Get Timeout with 200 status"
req=$URL/timeout
GET $req
check_ok_status "$out"

try "FC1 read bits test: port 502 - length 1 - (2/5)"
req="$URL/fc/1?ip=$SLAVE&port=502&slave=1&addr=8&len=1"
GET $req
check_ok_status "$out"

try "FC1 read bits test: port 502 - length 7 - (3/5)"
req="$URL/fc/1?ip=$SLAVE&port=502&slave=1&addr=3&len=7"
GET $req
check_ok_status "$out"

try "FC1 read bits test: port 502 - Illegal data address - (4/5)"
req="$URL/fc/1?ip=$SLAVE&port=502&slave=1&addr=20000&len=7"
GET $req
check_not_ok_status "$out"

try "FC1 read bits test: port 503 - length 7 - (5/5)"
req="$URL/fc/1?ip=$SLAVE&port=503&slave=1&addr=3&len=7"
GET $req
check_ok_status "$out"



# 5
#try "FC1 read bits test: port 503 - length 7 - (5/5)"
#out=$(curl -sL "$URL/fc/1?ip=$SLAVE&port=503&slave=1&addr=3&len=7")
#check_ok_status "$out"


#out=$(curl -s -w "%{http_code}" $URL)
#assert "404" "$out"

#try "Example of POST XML"

# Post xml (from hello.xml file) on /hello
#out=$(cat test/hello.xml | curl -s -H "Content-Type: text/xml" -d @- \
#  -X POST $URL/hello)
#assert "Hello World" "$out"

###############################################################
echo
echo "PASS: $tests_run tests run"
echo "FAIL: $fail_run tests run"

