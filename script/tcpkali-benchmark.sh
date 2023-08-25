#!/bin/bash
SEC="30s"
SLEEP_SEC="40"

echo "quickws.delay.12x:"
killall quickws.linux &>/dev/null
#./quickws.linux --use-delay-write --addr ":9000" &
./quickws.linux -b 12 --use-delay-write -a ":9000" --delay-write-init-buffer-size 11264 &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC
echo ""
echo "quickws.delay.8x:"
killall quickws.linux &>/dev/null

#./quickws.linux --use-delay-write --addr ":9000" &
./quickws.linux -b 8 --use-delay-write -a ":9000" --delay-write-init-buffer-size 11264 &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC
echo ""

echo "quickws.tcp-delay:"
killall quickws.linux &>/dev/null
#./quickws.linux --use-delay-write --addr ":9000" &
./quickws.linux -o -a ":9000"  &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC
echo ""

echo "quickws.delay.tcp-delay.windows.8x:"
killall quickws.linux &>/dev/null
./quickws.linux -o -b 8 --use-delay-write -a ":9000" --delay-write-init-buffer-size 11264 &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC
echo ""

echo "quickws.bufio:"
killall quickws.linux &>/dev/null
./quickws.linux -u --addr ":9000" &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC
echo ""

echo "quickws.bufio:8x"
killall quickws.linux &>/dev/null
./quickws.linux -u -b 8 --addr ":9000" &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC
echo ""

echo "quickws.windows.8x:"
killall quickws.linux &>/dev/null
./quickws.linux -w 8 --addr ":9000" &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC
echo ""

echo "quickws.windows.4x:"
killall quickws.linux &>/dev/null
./quickws.linux -w 4 --addr ":9000" &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC
echo ""

echo "quickws.windows:"
killall quickws.linux &>/dev/null
./quickws.linux --addr ":9000" &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep $SLEEP_SEC
echo ""

echo "gws:"
killall gws-std.linux &>/dev/null
./gws-std.linux --addr ":9001" &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9001/
kill $PID
sleep $SLEEP_SEC
echo ""

echo "gws.asyncwrite:"
killall gws-std.linux &>/dev/null
./gws-std.linux -a --addr ":9001" &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9001/
kill $PID
sleep $SLEEP_SEC
echo ""

echo "gws-std:"
killall gws-std.linux &>/dev/null
./gws-std.linux --addr ":9002" &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9002/
kill $PID
sleep $SLEEP_SEC
echo ""

echo "gws-std.asyncwrite:"
killall gws-std.linux &>/dev/null
./gws-std.linux -a --addr ":9002" &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9002/
kill $PID
sleep $SLEEP_SEC
echo ""

echo "gorilla-linux-ReadMessage:"
killall gorilla.linux  &>/dev/null
./gorilla.linux --addr ":9003" &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9003/
kill $PID
sleep $SLEEP_SEC
echo ""

echo "gorilla-linux-UseReader:"
killall gorilla.linux &>/dev/null
./gorilla.linux --addr ":9004" -u &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9004/
kill $PID
sleep $SLEEP_SEC
echo ""

echo "nettyws:"
killall nettyws.linux &>/dev/null
./nettyws.linux --addr ":9005" &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9005/
kill $PID
sleep $SLEEP_SEC
echo ""

echo "gobwas:"
killall gobwas &>/dev/null
./gobwas.linux --addr ":9006" &
PID=$!
sleep 1
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9006/
kill $PID
sleep $SLEEP_SEC
