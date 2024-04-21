#!/bin/bash

WS_PID=0
DEFAULT_DIR="$1"
if [ -z $DEFAULT_DIR ];then
    DEFAULT_DIR="./output/greatws-quickws"
fi

EXE="linux"
if [ `uname -s` == "Darwin" ];then
    EXE="mac"
fi

function ctrl_c_handler() {
    echo "Ctrl+C pressed. Exiting gracefully."
    if [ $WS_PID -eq 0 ];then
        exit 0
    fi 
    kill $WS_PID
    exit 0
}

trap ctrl_c_handler SIGINT

function tps_greatws_io() {
    # test tps
    echo "## TPS scenario testing"

    echo "#### greatws runs on I/O goroutines-greatws bind mode"
    pkill greatws 2>/dev/null
    ./greatws.linux -r &>/dev/null &
    WS_PID=$!
    sleep 1
    FILE_NAME="greatws-io-event"
    ./bench-ws.linux -c 10000 -d 60s -w "ws://127.0.0.1:23001/ws" --open-tmp-result --conns 10000 --JSON --label $FILE_NAME &> "$DEFAULT_DIR/$FILE_NAME.tmp.json"
    kill $WS_PID
}
function tps_greatws_stream2() {
    echo "#### greatws runs on business goroutines"
    pkill greatws 2>/dev/null
    sleep 1
    ./greatws.linux &>/dev/null &
    WS_PID=$!
    sleep 1
    FILE_NAME="greatws-execlist"
    ./bench-ws.linux -c 10000 -d 60s -w "ws://127.0.0.1:23001/ws" --open-tmp-result --conns 10000 --JSON --label $FILE_NAME &> "$DEFAULT_DIR/$FILE_NAME.tmp.json"
    kill $WS_PID
}

function tps_greatws_unstream() {
    echo "#### greatws runs on business goroutines(unstream)"
    pkill greatws 2>/dev/null
    sleep 1
    ./greatws.linux -u &>/dev/null &
    WS_PID=$!
    sleep 1
    FILE_NAME="greatws-unstream"
    ./bench-ws.linux -c 10000 -d 60s -w "ws://127.0.0.1:23001/ws" --open-tmp-result --conns 10000 --JSON --label $FILE_NAME &> "$DEFAULT_DIR/$FILE_NAME.tmp.json"
    kill $WS_PID
}

function tps_greatws_stream() {
    echo "#### greatws uses one Goroutine per connection"
    pkill greatws 2>/dev/null
    sleep 1
    ./greatws.linux -s &>/dev/null &
    WS_PID=$!
    sleep 1
    FILE_NAME="greatws-stream"
    ./bench-ws.linux -c 10000 -d 60s -w "ws://127.0.0.1:23001/ws" --open-tmp-result --conns 10000 --JSON --label $FILE_NAME &> "$DEFAULT_DIR/$FILE_NAME.tmp.json"
    kill $WS_PID
}

function tps_quickws() {
    echo "#### quickws"
    pkill quickws 2>/dev/null
    sleep 1
    ./quickws.linux &>/dev/null &
    WS_PID=$!
    sleep 1
    FILE_NAME="quickws"
    ./bench-ws.linux -c 10000 -d 60s -w "ws://127.0.0.1:23001/ws" --open-tmp-result --conns 10000 --JSON --label $FILE_NAME &> "$DEFAULT_DIR/$FILE_NAME.tmp.json"
    kill $WS_PID
}

function tps_quickws_mini() {
    pkill quickws 2>/dev/null
    sleep 1
    ./bin/"quickws.$EXE" -l -1 &>/dev/null &
    WS_PID=$!
    sleep 1
    FILE_NAME="quickws_mini"
    #./bin/"bench-ws.$EXE" -c 100 -d 10s -w "ws://127.0.0.1:23001/ws" --conns 100 --JSON --label "$FILE_NAME"  --open-tmp-result
    ./bin/"bench-ws.$EXE" -c 100 -d 10s -w "ws://127.0.0.1:23001/ws" --conns 100 --JSON --label "$FILE_NAME" &> "$DEFAULT_DIR/$FILE_NAME.tmp.json"
    kill $WS_PID
}

function tps_test_mini() {
    if [ ! -d "$DEFAULT_DIR" ];then
        mkdir -p $DEFAULT_DIR
    fi

    tps_quickws_mini
}

function tps_test() {
    if [ ! -d "$DEFAULT_DIR" ];then
        mkdir -p $DEFAULT_DIR
    fi

    tps_greatws_io
    tps_greatws_stream2
    tps_greatws_unstream
    tps_greatws_stream
    tps_quickws
}

function traffic_test() {
    # test 流量
    echo "## traffic scenario testing"
    echo "#### greatws runs on business goroutines"
    pkill greatws 2>/dev/null
    sleep 1 #防止进程还存在
    ./greatws.linux &>/dev/null &
    WS_PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f ./testdata/1K.txt --ws 127.0.0.1:9001/
    kill $WS_PID

    echo "#### greatws runs on I/O goroutines"
    pkill greatws 2>/dev/null
    sleep 1 #防止进程还存在
    ./greatws.linux -r &>/dev/null &
    WS_PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f ./testdata/1K.txt --ws 127.0.0.1:9001/
    kill $WS_PID

    echo "#### greatws runs on business goroutines unstream"
    pkill greatws 2>/dev/null
    sleep 1 #防止进程还存在
    ./greatws.linux -u &>/dev/null &
    WS_PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f ./testdata/1K.txt --ws 127.0.0.1:9001/
    kill $WS_PID

    echo "#### greatws uses one Goroutine per connection"
    pkill greatws 2>/dev/null
    sleep 1 #防止进程还存在
    ./greatws.linux -s &>/dev/null &
    WS_PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f ./testdata/1K.txt --ws 127.0.0.1:9001/
    kill $WS_PID

    echo "#### quickws"
    pkill greatws 2>/dev/null
    pkill quickws 2>/dev/null
    sleep 1 #防止进程还存在
    ./quickws.linux &>/dev/null &
    WS_PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f ./testdata/1K.txt --ws 127.0.0.1:9001/
    kill $WS_PID
}

function build() {
    echo "## build"
    make
}

tps_test_mini
#build
#tps_test
#traffic_test

