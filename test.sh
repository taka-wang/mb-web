#!/bin/bash

#################################################
# REST API Tester
#
# Author: Taka Wang
# Date: 2016/09/13
#################################################

## Varaibles
URL=http://localhost:8080/api/mb/tcp
SLAVE=$(docker inspect --format '{{ .NetworkSettings.IPAddress }}' mbweb_slave_1)

# color code ---------------
COLOR_REST='\e[0m'
COLOR_GREEN='\e[1;32m';
COLOR_RED='\e[1;31m';

## Unit-Testable Shell Scripts (http://eradman.com/posts/ut-shell-scripts.html)
typeset -i tests_run=0
typeset -i fail_run=0
function try { 
    this="$1"
    if [ "$(uname)" == "Darwin" ]; then
        echo "### " $this    
    elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
        echo -e "${COLOR_RED}### $this${COLOR_REST}"
    fi
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
function set_title {
    if [ "$(uname)" == "Darwin" ]; then
        echo "========== $1 =========="
    elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
        echo -e "${COLOR_GREEN}========== $1 ==========${COLOR_REST}"
    fi
}

# HTTP GET
function GET {
    echo "req :" $1
    out=$(curl -sL $1)
}
## end

###############################################################
echo "---------------------------------" # start

#######################################################
set_title "TestTimeoutOps"
#######################################################

try "Get timeout 1st round - (1/4)"
req=$URL/timeout
GET "$req"
check_ok_status "$out"

try "Set valid timeout - (2/4)"
out=$(curl -sL -X POST -H "Content-Type: application/json"  -d '{ "timeout": 210000 }' "$URL/timeout")
check_ok_status "$out"

try "Get timeout 2nd round - (3/4)"
req=$URL/timeout
GET "$req"
check_ok_status "$out"

try "Set invalid timeout - (4/4)"
out=$(curl -sL -X POST -H "Content-Type: application/json"  -d '{ "timeout": 123 }' "$URL/timeout")
check_not_ok_status "$out"

#######################################################
set_title "TestOneOffWriteFC5"
#######################################################

try "FC5 write bit test: port 502 - invalid value(2) - (1/5)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "503",
    "slave": 1,
    "addr": 10,
    "data": 2
    }' "$URL/fc/5")
check_ok_status "$out"

try "FC5 write bit test: port 502 - miss port - (2/5)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "slave": 1,
    "addr": 10,
    "data": 2
    }' "$URL/fc/5")
check_ok_status "$out"

try "FC5 write bit test: port 502 - valid value(0) - (3/5)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "data": 0
    }' "$URL/fc/5")
check_ok_status "$out"

try "FC5 write bit test: port 502 - valid value(1) - (4/5)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "data": 1
    }' "$URL/fc/5")
check_ok_status "$out"

try "FC5 invalid function code - (5/5)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "data": 1
    }' "$URL/fc/7")
check_not_ok_status "$out"

#######################################################
set_title "TestOneOffWriteFC6"
#######################################################

try "FC6 write 'DEC' register test: port 502 - valid value (22) - (1/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "hex": false,
	"data": "22"
    }' "$URL/fc/6")
check_ok_status "$out"

try "FC6 write 'DEC' register test: port 502 - miss hex type & port - (2/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "slave": 1,
	"data": "22"
    }' "$URL/fc/6")
check_ok_status "$out"

try "FC6 write 'DEC' register test: port 502 - invalid value (array) - (3/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "hex": false,
    "data": "22,11"
    }' "$URL/fc/6")
check_ok_status "$out"

try "FC6 write 'DEC' register test: port 502 - invalid hex type - (4/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "hex": false,
    "data": "ABCD1234"
    }' "$URL/fc/6")
check_not_ok_status "$out"

try "FC6 write 'HEX' register test: port 502 - valid value (ABCD) - (5/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "hex": true,
    "data": "ABCD"
    }' "$URL/fc/6")
check_ok_status "$out"

try "FC6 write 'HEX' register test: port 502 - miss port (ABCD) - (6/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "slave": 1,
    "addr": 10,
    "hex": true,
    "data": "ABCD"
    }' "$URL/fc/6")
check_ok_status "$out"

try "FC6 write 'HEX' register test: port 502 - invalid value (ABCD1234) - (7/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "hex": true,
    "data": "ABCD1234"
    }' "$URL/fc/6")
check_ok_status "$out"

try "FC6 write 'HEX' register test: port 502 - invalid hex type - (8/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "hex": true,
    "data": "22,11"
    }' "$URL/fc/6")
check_not_ok_status "$out"

#######################################################
set_title "TestOneOffWriteFC15"
#######################################################

try "FC15 write bit test: port 502 - invalid json type - (1/5)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "data": [-1,0,-1,0]
    }' "$URL/fc/15")
check_not_ok_status "$out"

try "FC15 write bit test: port 502 - invalid json type - (2/5)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "data": "1,0,1,0"
    }' "$URL/fc/15")
check_not_ok_status "$out"

try "FC15 write bit test: port 502 - invalid value(2) - (3/5)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "data": [2,0,2,0]
    }' "$URL/fc/15")
check_ok_status "$out"

try "FC15 write bit test: port 502 - miss from & port - (4/5)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "data": [2,0,2,0]
    }' "$URL/fc/15")
check_ok_status "$out"

try "FC15 write bit test: port 502 - valid value(0) - (5/5)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "data": [0,1,1,0]
    }' "$URL/fc/15")
check_ok_status "$out"

#######################################################
set_title "TestOneOffWriteFC16"
#######################################################

try "FC16 write write 'DEC' register test: port 502 - valid value (11,22,33,44) - (1/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "hex": false,
    "data": "11,22,33,44"
    }' "$URL/fc/16")
check_ok_status "$out"

try "FC16 write write 'DEC' register test: port 502 - miss hex type & port - (2/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "data": "11,22,33,44"
    }' "$URL/fc/16")
check_ok_status "$out"

try "FC16 write write 'DEC' register test: port 502 - invalid hex type - (3/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "hex": false,
    "data": "ABCD1234"
    }' "$URL/fc/16")
check_not_ok_status "$out"

try "FC16 write write 'DEC' register test: port 502 - invalid length - (4/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "len": 8,
    "hex": false,
    "data": "11,22,33,44"
    }' "$URL/fc/16")
check_ok_status "$out"

try "FC16 write write 'HEX' register test: port 502 - valid value (ABCD1234) - (5/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "hex": true,
    "data": "ABCD1234"
    }' "$URL/fc/16")
check_ok_status "$out"

try "FC16 write write 'HEX' register test: port 502 - miss port (ABCD) - (6/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "hex": true,
    "data": "ABCD1234"
    }' "$URL/fc/16")
check_ok_status "$out"

try "FC16 write write 'HEX' register test: port 502 - invalid hex type (11,22,33,44) - (7/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "len": 4,
    "hex": true,
    "data": "11,22,33,44"
    }' "$URL/fc/16")
check_not_ok_status "$out"

try "FC16 write write 'HEX' register test: port 502 - invalid length - (8/8)"
out=$(curl -sL -H "Content-Type: application/json" -X POST  --data '{
    "ip": "'"$SLAVE"'",
    "port": "502",
    "slave": 1,
    "addr": 10,
    "len": 8,
    "hex": true,
    "data": "ABCD1234"
    }' "$URL/fc/16")
check_ok_status "$out"

#######################################################
set_title "TestOneOffReadFC1"
#######################################################

try "FC1 read bits test: port 502 - miss ip - (1/5)"
req="$URL/fc/1?&port=502&slave=1&addr=3&len=1"
GET "$req"
check_not_ok_status "$out"

try "FC1 read bits test: port 502 - length 1 - (2/5)"
req="$URL/fc/1?ip=$SLAVE&port=502&slave=1&addr=8&len=1"
GET "$req"
check_ok_status "$out"

try "FC1 read bits test: port 502 - length 7 - (3/5)"
req="$URL/fc/1?ip=$SLAVE&port=502&slave=1&addr=3&len=7"
GET "$req"
check_ok_status "$out"

try "FC1 read bits test: port 502 - Illegal data address - (4/5)"
req="$URL/fc/1?ip=$SLAVE&port=502&slave=1&addr=20000&len=7"
GET "$req"
check_not_ok_status "$out"

try "FC1 read bits test: port 503 - length 7 - (5/5)"
req="$URL/fc/1?ip=$SLAVE&port=503&slave=1&addr=3&len=7"
GET "$req"
check_ok_status "$out"

#######################################################
set_title "TestOneOffReadFC2"
#######################################################

try "FC2 read bits test: port 502 - length 1 - (1/4)"
req="$URL/fc/2?ip=$SLAVE&port=502&slave=1&addr=3&len=1"
GET "$req"
check_ok_status "$out"

try "FC2 read bits test: port 502 - length 7 - (2/4)"
req="$URL/fc/2?ip=$SLAVE&port=502&slave=1&addr=3&len=7"
GET "$req"
check_ok_status "$out"

try "FC2 read bits test: port 502 - Illegal data address - (3/4)"
req="$URL/fc/2?ip=$SLAVE&port=502&slave=1&addr=20000&len=7"
GET "$req"
check_not_ok_status "$out"

try "FC2 read bits test: port 503 - length 7 - (4/4)"
req="$URL/fc/2?ip=$SLAVE&port=503&slave=1&addr=3&len=7"
GET "$req"
check_ok_status "$out"

#######################################################
set_title "TestOneOffReadFC3"
#######################################################

try "FC3 read bytes Type 1 test: port 502 - (1/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=7&type=1"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 2 test: port 502 - (2/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=7&type=2"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 3 length 4 test: port 502 - (3/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=4&type=3&a=0&b=65535&c=100&d=500"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 3 length 7 test: port 502 - invalid length - (4/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=7&type=3&a=0&b=65535&c=100&d=500"
GET "$req"
check_not_ok_status "$out"

try "FC3 read bytes Type 4 length 4 test: port 502 - Order: AB - (5/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=4&type=4&order=1"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 4 length 4 test: port 502 - Order: BA - (6/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=4&type=4&order=2"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 4 length 4 test: port 502 - miss order - (7/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=4&type=4"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 5 length 4 test: port 502 - Order: AB - (8/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=4&type=5&order=1"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 5 length 4 test: port 502 - Order: BA - (9/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=4&type=5&order=2"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 5 length 4 test: port 502 - miss order - (10/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=4&type=5"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 6 length 8 test: port 502 - Order: AB - (11/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=6&order=1"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 6 length 8 test: port 502 - Order: BA - (12/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=6&order=2"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 6 length 8 test: port 502 - miss order - (13/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=6"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 6 length 7 test: port 502 - invalid length - (14/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=7&type=6&order=2"
GET "$req"
check_not_ok_status "$out"

try "FC3 read bytes Type 7 length 8 test: port 502 - Order: AB - (15/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=7&order=1"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 7 length 8 test: port 502 - Order: BA - (16/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=7&order=2"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 7 length 8 test: port 502 - miss order - (17/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=7"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 7 length 7 test: port 502 - invalid length - (18/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=7&type=7&order=2"
GET "$req"
check_not_ok_status "$out"

try "FC3 read bytes Type 8 length 8 test: port 502 - order: ABCD - (19/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=8&order=1"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 8 length 8 test: port 502 - order: DCBA - (20/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=8&order=2"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 8 length 8 test: port 502 - order: BADC - (21/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=8&order=3"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 8 length 8 test: port 502 - order: CDAB - (22/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=8&order=4"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 8 length 7 test: port 502 - invalid length - (23/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=7&type=8&order=1"
GET "$req"
check_not_ok_status "$out"

try "FC3 read bytes: port 502 - invalid type - (24/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=9&order=1"
GET "$req"
check_ok_status "$out"

#######################################################
set_title "TestOneOffReadFC4"
#######################################################

try "FC4 read bytes Type 1 test: port 502 - (1/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=7&type=1"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 2 test: port 502 - (2/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=7&type=2"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 3 length 4 test: port 502 - (3/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=4&type=3&a=0&b=65535&c=100&d=500"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 3 length 7 test: port 502 - invalid length - (4/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=7&type=3&a=0&b=65535&c=100&d=500"
GET "$req"
check_not_ok_status "$out"

try "FC4 read bytes Type 4 length 4 test: port 502 - Order: AB - (5/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=4&type=4&order=1"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 4 length 4 test: port 502 - Order: BA - (6/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=4&type=4&order=2"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 4 length 4 test: port 502 - miss order - (7/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=4&type=4"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 5 length 4 test: port 502 - Order: AB - (8/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=4&type=5&order=1"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 5 length 4 test: port 502 - Order: BA - (9/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=4&type=5&order=2"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 5 length 4 test: port 502 - miss order - (10/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=4&type=5"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 6 length 8 test: port 502 - Order: AB - (11/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=6&order=1"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 6 length 8 test: port 502 - Order: BA - (12/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=6&order=2"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 6 length 8 test: port 502 - miss order - (13/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=6"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 6 length 7 test: port 502 - invalid length - (14/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=7&type=6&order=2"
GET "$req"
check_not_ok_status "$out"

try "FC4 read bytes Type 7 length 8 test: port 502 - Order: AB - (15/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=7&order=1"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 7 length 8 test: port 502 - Order: BA - (16/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=7&order=2"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 7 length 8 test: port 502 - miss order - (17/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=7"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 7 length 7 test: port 502 - invalid length - (18/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=7&type=7&order=2"
GET "$req"
check_not_ok_status "$out"

try "FC4 read bytes Type 8 length 8 test: port 502 - order: ABCD - (19/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=8&order=1"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 8 length 8 test: port 502 - order: DCBA - (20/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=8&order=2"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 8 length 8 test: port 502 - order: BADC - (21/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=8&order=3"
GET "$req"
check_ok_status "$out"

try "FC4 read bytes Type 8 length 8 test: port 502 - order: CDAB - (22/24)"
req="$URL/fc/3?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=8&order=4"
GET "$req"
check_ok_status "$out"

try "FC3 read bytes Type 8 length 7 test: port 502 - invalid length - (23/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=7&type=8&order=1"
GET "$req"
check_not_ok_status "$out"

try "FC4 read bytes: port 502 - invalid type - (24/24)"
req="$URL/fc/4?ip=$SLAVE&port=502&slave=1&addr=3&len=8&type=9&order=1"
GET "$req"
check_ok_status "$out"

###############################################################
echo
if [ "$(uname)" == "Darwin" ]; then
    echo "PASS: $tests_run tests run"
    echo "FAIL: $fail_run tests run"
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    echo -e "${COLOR_GREEN}PASS: $tests_run tests run${COLOR_REST}"
    echo -e "${COLOR_RED}FAIL: $fail_run tests run${COLOR_REST}"
fi
echo "---------------------------------" # end



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