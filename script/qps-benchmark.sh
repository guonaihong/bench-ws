#!/bin/bash
SEC="18s"
SLEEP_SEC="40"

echo "quickws.1:"
killall quickws.linux &>/dev/null
./quickws.linux --addr ":9000" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9000/" -c 10000 -d $SEC
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC

echo "quickws.2:"
killall quickws.linux &>/dev/null
./quickws.linux --addr ":9000" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9000/ws" -c 10000 -d $SEC
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC

echo "gws:"
killall gws-std.linux &>/dev/null
./gws-std.linux --addr ":9001" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9001/ws" -c 10000 -d $SEC
kill $PID
sleep $SLEEP_SEC

echo "gws-std:"
killall gws-std.linux &>/dev/null
./gws-std.linux --addr ":9002" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9002/ws" -c 10000 -d $SEC
kill $PID
sleep $SLEEP_SEC

echo "gorilla-linux-ReadMessage:"
killall gorilla.linux  &>/dev/null
./gorilla.linux --addr ":9003" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9003/ws" -c 10000 -d $SEC
kill $PID
sleep $SLEEP_SEC

echo "gorilla-linux-UseReader:"
killall gorilla.linux &>/dev/null
./gorilla.linux --addr ":9004" -u &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9004/ws" -c 10000 -d $SEC
kill $PID
sleep $SLEEP_SEC

echo "nettyws:"
killall nettyws.linux &>/dev/null
./nettyws.linux --addr ":9005" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9005/ws" -c 10000 -d $SEC
kill $PID
sleep $SLEEP_SEC

echo "gobwas:"
killall gobwas &>/dev/null
./gobwas.linux --addr ":9006" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9006/ws" -c 10000 -d $SEC
kill $PID
sleep $SLEEP_SEC
