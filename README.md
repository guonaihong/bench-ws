# my-bench-ws
压测websocket(quickws)的仓库

# 压测
* 启动服务端的命令
```
./autobahn-quickws.linux -d
```
* 启动客户端的命令
```
./test-client.linux --close-check -c 10000 -t 1000000000 -w "ws://127.0.0.1:9001/autobahn"
```
