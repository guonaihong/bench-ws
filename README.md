# bench-ws

压测websocket(quickws)的仓库

### 一、如何在界面观看柱状图
```console
./script/start-web.sh
# 找到App running at: 这一行, 浏览器点开直接查看
```
### 二、单服务跑tps压测

* 启动服务端的命令

```
./bin/quickws.linux
```

* 启动客户端的命令

```
./bin/bench-ws.linux --close-check -c 10000 -t 1000000000 -w "ws://127.0.0.1:9001/"
```

### 三、单服务跑流量压测

* 启动服务端的命令

```
./bin/quickws.linux -o -u
```

* 启动客户端的命令

```
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9001/
```

### 四、批量跑tps压测

```
make
./script/tps-all-benchmark.sh
```