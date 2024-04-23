#!/bin/bash

# Define variables
SEC="18s"
SLEEP_SEC="1"
TEST_SEC="10s"
BIN="./bin"
EXE="linux"
DEFAULT_DIR="./output/all"

if [ ! -d "$DEFAULT_DIR" ]; then
        mkdir -p $DEFAULT_DIR
fi

# Determine the operating system type and set the EXE variable accordingly
if [ "$(uname -s)" == "Darwin" ]; then
    EXE="mac"
elif [ "$(uname -s)" == "FreeBSD" ]; then
    EXE="freebsd"
fi

# Function to kill all running WebSocket servers
function kill_all_servers() {
    echo "Killing all WebSocket servers..."
    killall greatws.linux quickws.linux gws-std.linux \
            gorilla.linux nettyws.linux gobwas.linux \
            nbio-std.linux nbio-nonblocking.linux \
            nbio-blocking.linux nbio-mixed.linux \
            hertz.linux hertz-std.linux \
            fasthttp-ws-std.linux &>/dev/null
    echo "All WebSocket servers killed."
}

# Function to run test for a WebSocket server
function run_test() {
    local server_name="$1"
    local address="$2"
    local suffix="$3"
    local ws_arg="$"

    echo "$server_name:"
    kill_all_servers
    "$BIN/$server_name.$EXE" $ws_arg --addr ":$address" &>/dev/null &
    local PID=$!
    sleep 1
    FILE_NAME="${server_name}_${suffix}"
    "$BIN/bench-ws.$EXE" -c 10000 -d $TEST_SEC -w "ws://127.0.0.1:$address/ws" \
                         --conns 10000 --JSON --label "$FILE_NAME" \
                         &> "$DEFAULT_DIR/$FILE_NAME.tmp.json"
    kill "$PID"
    sleep "$SLEEP_SEC"
}

# Run tests for each WebSocket server
run_test "greatws" "24001" "1"
run_test "quickws" "23001" "0"
run_test "quickws" "23001" "1"
run_test "gws"  "13001" "2"
run_test "gws-std" "14001" "3"
run_test "gorilla" "12001" "4.1" "-u"
run_test "gorilla" "12001" "4.2"
run_test "nettyws" "21001" "5"
run_test "gobwas" "11001" "6"
run_test "nbio-std" "20001" "7"
run_test "nbio-nonblocking" "19001" "8"
run_test "nbio-blocking" "17001" "9"
run_test "nbio-mixed" "18001" "10"
run_test "hertz" "15001" "11"
run_test "hertz-std" "16001" "12"
run_test "fasthttp-ws-std" "10001" "13"

./script/release.sh