#!/bin/bash

WS_PID=0
DEFAULT_DIR="$1"
BIN="./bin"
TEST_SEC="10s"
if [ -z $DEFAULT_DIR ]; then
    DEFAULT_DIR="./output/greatws-quickws"
fi

EXE="linux"
MAKE_TARGET="build-linux"

if [ `uname -s` == "Darwin" ]; then
    EXE="mac"
    MAKE_TARGET="build-mac"
fi

if [ `uname -s` == "MINGW"* ]; then
    EXE="windows"
    MAKE_TARGET="build-windows"
fi

function ctrl_c_handler() {
    echo "Ctrl+C pressed. Exiting gracefully."
    if [ $WS_PID -ne 0 ]; then
        kill $WS_PID
    fi 
    exit 0
}

function kill_server() {
    pkill greatws 2>/dev/null
    pkill quickws 2>/dev/null
}

trap ctrl_c_handler SIGINT

function build_executables() {
    echo "## Building executables for $EXE"
    make $MAKE_TARGET
}

function run_test() {
    local BIN_NAME="$1"
    local WS_ARGS="$2"
    local FILE_SUFFIX="$3"
    local address="$4"

    kill_server
    sleep 1
    $BIN/$BIN_NAME.$EXE $WS_ARGS &>/dev/null &
    WS_PID=$!
    sleep 1
    FILE_NAME="$BIN_NAME-$FILE_SUFFIX"
    $BIN/bench-ws.$EXE -c 10000 -d $TEST_SEC -w "ws://127.0.0.1:$address/ws" --conns 10000 --JSON --label $FILE_NAME &> "$DEFAULT_DIR/$FILE_NAME.tmp.json"
    kill $WS_PID
}

function tps_greatws_io() {
    echo "## TPS scenario testing"
    echo "#### greatws runs on I/O goroutines-greatws bind mode"
    run_test "greatws" "-r" "io-event" "24001"
}

function tps_greatws_stream2() {
    echo "#### greatws runs on business goroutines"
    run_test "greatws" "execlist" "24001"
}

function tps_greatws_unstream() {
    echo "#### greatws runs on business goroutines(unstream)"
    run_test "greatws" "-u" "unstream" "24001"
}

function tps_greatws_stream() {
    echo "#### greatws uses one Goroutine per connection"
    run_test "greatws" "-s" "stream" "24001"
}

function tps_quickws() {
    echo "#### quickws"
    run_test "quickws" "" "quickws" "23001"
}

function tps_quickws_mini() {
    kill_server
    sleep 1
    $BIN/quickws.$EXE -l -1 &>/dev/null &
    WS_PID=$!
    sleep 1
    FILE_NAME="quickws_mini"
    $BIN/bench-ws.$EXE -c 100 -d 10s -w "ws://127.0.0.1:23001/ws" --conns 100 --JSON --label "$FILE_NAME" &> "$DEFAULT_DIR/$FILE_NAME.tmp.json"
    kill $WS_PID
}

function tps_test_debug() {
    if [ ! -d "$DEFAULT_DIR" ]; then
        mkdir -p $DEFAULT_DIR
    fi

    tps_greatws_stream2
}

function tps_test() {
    if [ ! -d "$DEFAULT_DIR" ]; then
        mkdir -p $DEFAULT_DIR
    fi

    build_executables

    tps_greatws_io
    tps_greatws_stream2
    tps_greatws_unstream
    tps_greatws_stream
    tps_quickws
}

function traffic_test() {
    echo "## traffic scenario testing"

    for WS_TYPE in "" "-r" "-u" "-s"; do
        echo "#### greatws runs on business goroutines$([[ ! -z $WS_TYPE ]] && echo "($WS_TYPE)")"
        kill_server
        sleep 1
        $BIN/greatws.$EXE $WS_TYPE &>/dev/null &
        WS_PID=$!
        sleep 1
        tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f ./testdata/1K.txt --ws 127.0.0.1:9001/
        kill $WS_PID
    done

    echo "#### quickws"
    kill_server
    sleep 1
    $BIN/quickws.$EXE &>/dev/null &
    WS_PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f ./testdata/1K.txt --ws 127.0.0.1:9001/
    kill $WS_PID
}

function build() {
    echo "## build"
    make
}

tps_test
