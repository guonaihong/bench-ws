## 编译

```console
make
```

## 一、不限制cpu，多端口测试

* 1.启动服务端

```console
./quickws.linux  -l -1
```

* 2.启动客户端

## 二、限制cpu，多端口测试

这里用的5800h，8c16t，所以0-7号逻辑cpu绑定服务端,8-15号逻辑cpu绑定客户端

* 2.1 启动服务端

```console
taskset -c 0-7 ./quickws.linux  -l -1
```

* 2.2 启动客户端

```console
./bench-ws.linux -c 10000 -w "ws://127.0.0.1:23001-23050/ws" -c 10000 -d 10s --open-tmp-result
```
