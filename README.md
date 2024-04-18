# bench-ws

压测websocket(quickws)的仓库

# tps压测

* 启动服务端的命令

```
./bin/quickws.linux
```

* 启动客户端的命令

```
./bin/test-client.linux --close-check -c 10000 -t 1000000000 -w "ws://127.0.0.1:9001/"
```

# 流量压测

* 启动服务端的命令

```
./bin/quickws.linux -o -u
```

* 启动客户端的命令

```
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9001/
```

# 批量跑tps压测

```
make
./script/tps-benchmark.sh
```
