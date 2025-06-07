# bench-ws

压测websocket(quickws)的仓库

### 一、终端压测
```
./benchmark/benchmark-c10k.sh
```

### 一、如何在界面观看柱状图

#### 1.1 依赖安装

安装nvm（Node Version Manager）：  
nvm是一个Node.js版本管理工具，它允许你安装和使用不同版本的Node.js和npm。

安装nvm：

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
```

重新加载你的终端或运行以下命令来使nvm生效：

```bash
source ~/.nvm/nvm.sh
```

安装新版本的Node.js和npm：
使用nvm安装最新版本的Node.js和npm：

```bash
nvm install node  # 安装最新稳定版本的Node.js和npm
```

切换到新版本的Node.js：

`
nvm use node
`
安装yarn

```bash
npm install -g yarn
```

#### 1.2 运行web服务

```console
./script/start-web.sh
# 找到App running at: 这一行, 浏览器点开直接查看
```

### 二、单服务跑tps压测

* 启动服务端的命令

```bash
./bin/quickws.linux
```

* 启动客户端的命令

```bash
./bin/bench-ws.linux --close-check -c 10000 -t 1000000000 -w "ws://127.0.0.1:9001/"
```

### 三、单服务跑流量压测

* 启动服务端的命令

```bash
./bin/quickws.linux -o -u
```

* 启动客户端的命令

```bash
tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9001/
```


