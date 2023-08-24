#!/bin/bash

# 小主机散热一般，可以用双倍时间散下热
SEC="18s"
SLEEP_SEC="40"

echo "quickws.0:忽略第一个的成绩"
killall quickws.linux &>/dev/null
./quickws.linux --addr ":9000" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9000/" -c 10000 -d $SEC
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC

echo "quickws.1:"
killall quickws.linux &>/dev/null
./quickws.linux --addr ":9000" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9000/ws" -c 10000 -d $SEC
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC

echo "gws.2:"
killall gws-std.linux &>/dev/null
./gws.linux --addr ":9001" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9001/ws" -c 10000 -d $SEC
kill $PID
sleep $SLEEP_SEC

echo "gws-std.3:"
killall gws-std.linux &>/dev/null
./gws-std.linux --addr ":9002" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9002/ws" -c 10000 -d $SEC
kill $PID
sleep $SLEEP_SEC

echo "gorilla-linux-ReadMessage.4.1:"
killall gorilla.linux  &>/dev/null
./gorilla.linux --addr ":9003" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9003/ws" -c 10000 -d $SEC
kill $PID
sleep $SLEEP_SEC

echo "gorilla-linux-UseReader.4.2:"
killall gorilla.linux &>/dev/null
./gorilla.linux --addr ":9004" -u &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9004/ws" -c 10000 -d $SEC
kill $PID
sleep $SLEEP_SEC

echo "nettyws.5:"
killall nettyws.linux &>/dev/null
./nettyws.linux --addr ":9005" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9005/ws" -c 10000 -d $SEC
kill $PID
sleep $SLEEP_SEC

echo "gobwas.6:"
killall gobwas &>/dev/null
./gobwas.linux --addr ":9006" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9006/ws" -c 10000 -d $SEC
kill $PID
sleep $SLEEP_SEC

echo "nbio-std.7:"
killall quickws.linux &>/dev/null
./nbio-std.linux --addr ":9007" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9007/" -c 10000 -d $SEC
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC


echo "nbio-nonblocking.8:"
killall nbio-nonblocking.linux &>/dev/null
./nbio-nonblocking.linux --addr ":9008" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9008/" -c 10000 -d $SEC
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC


echo "nbio-blocking.9:"
killall nbio-blocking.linux &>/dev/null
./nbio-blocking.linux --addr ":9009" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9009/" -c 10000 -d $SEC
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC

echo "nbio-mixed.10:"
killall quickws.linux &>/dev/null
./nbio-mixed.linux --addr ":9010" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9010/" -c 10000 -d $SEC
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC

echo "hertz.11:"
killall hertz.linux &>/dev/null
./hertz.linux --addr ":9011" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9011/" -c 10000 -d $SEC
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC

echo "hertz-std.12:"
killall hertz-std.linux &>/dev/null
./hertz-std.linux --addr ":9012" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9012/" -c 10000 -d $SEC
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC


echo "fasthttp-ws.13:"
killall fasthttp-ws-std.linux &>/dev/null
./fasthttp-ws-std.linux  --addr ":9013" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9013/" -c 10000 -d $SEC
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC

