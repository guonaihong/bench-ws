echo "quickws:"
./quickws.linux --addr ":9000" &
PID=$!

sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9000/ws" -c 10000 -d 10s
kill $PID
# 让cpu 温度降些，防止过热，影响后观框架的测试, 散热好的, sleep时间可以改短些
sleep 15

echo "gws:"
./gws-std.linux --addr ":9001" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9001/ws" -c 10000 -d 10s
kill $PID
sleep 15

echo "gws-std:"
./gws-std.linux --addr ":9002" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9002/ws" -c 10000 -d 10s
kill $PID
sleep 15

echo "gorilla-linux-ReadMessage:"
./gorilla.linux --addr ":9003" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9003/ws" -c 10000 -d 10s
kill $PID
sleep 15

echo "gorilla-linux-UseReader:"
./gorilla.linux --addr ":9004" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9004/ws" -c 10000 -d 10s
kill $PID
sleep 15

echo "nettyws:"
./nettyws.linux --addr ":9005" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9005/ws" -c 10000 -d 10s
kill $PID
sleep 15

echo "gobwas:"
./gobwas.linux --addr ":9006" &
PID=$!
sleep 1
./test-client.linux -c 10000 -w "ws://127.0.0.1:9006/ws" -c 10000 -d 10s
kill $PID
sleep 15
