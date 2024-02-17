#!/bin/bash

WS_PID=0
function ctrl_c_handler() {
    echo "Ctrl+C pressed. Exiting gracefully."
    if [ $WS_PID -eq 0 ];then
        exit 0
    fi 
    kill $WS_PID
    exit 0
}

trap ctrl_c_handler SIGINT

function tps_test() {
    # test tps
    echo "## TPS scenario testing"

    echo "#### greatws runs on I/O goroutines-greatws bind mode"
    pkill greatws 2>/dev/null
    ./greatws.linux -r &>/dev/null &
    WS_PID=$!
    sleep 1
    ./test-client.linux -c 10000 -d 60s -w "ws://127.0.0.1:9001/autobahn" --open-tmp-result --conns 10000
    kill $WS_PID

    echo "#### greatws runs on business goroutines"
    pkill greatws 2>/dev/null
    sleep 1
    ./greatws.linux &>/dev/null &
    WS_PID=$!
    sleep 1
    ./test-client.linux -c 10000 -d 60s -w "ws://127.0.0.1:9001/autobahn" --open-tmp-result --conns 10000
    kill $WS_PID

    echo "#### greatws uses one Goroutine per connection"
    pkill greatws 2>/dev/null
    ./greatws.linux -s &>/dev/null &
    WS_PID=$!
    sleep 1
    ./test-client.linux -c 10000 -d 60s -w "ws://127.0.0.1:9001/autobahn" --open-tmp-result --conns 10000
    kill $WS_PID

    echo "#### quickws"
    pkill quickws 2>/dev/null
    ./quickws.linux &>/dev/null &
    WS_PID=$!
    sleep 1
    ./test-client.linux -c 10000 -d 60s -w "ws://127.0.0.1:9001/autobahn" --open-tmp-result --conns 10000
    kill $WS_PID
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

    echo "#### greatws runs on I/O goroutines-greatws bind mode"
    pkill greatws 2>/dev/null
    sleep 1 #防止进程还存在
    ./greatws.linux -r &>/dev/null &
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

tps_test
traffic_test

