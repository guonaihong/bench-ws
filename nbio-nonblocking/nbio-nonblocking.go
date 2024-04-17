package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	"github.com/guonaihong/bench-ws/config"
	"github.com/guonaihong/bench-ws/core"
	"github.com/guonaihong/clop"

	// "github.com/lesismal/nbio/log"
	"github.com/lesismal/nbio/mempool"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

var upgrader = websocket.NewUpgrader()

type Config struct {
	// 打开性能优化开关
	UseStdMalloc   bool `clop:"short;long" usage:"use reader"`
	ReadBufferSize int  `clop:"short;long" usage:"read buffer size" default:"1024"`

	Addr string `clop:"short;long" usage:"websocket server address" default:":4444""`
	core.BaseCmd
}

func main() {
	var cnf Config
	clop.Bind(&cnf)

	// 内存限制得越低效率越低、压测的带宽越低
	debug.SetMemoryLimit(1024 * 1024 * 512)
	if cnf.UseStdMalloc {
		// tcpkali这种场景，nbio用标准库比mempool内存占用低
		mempool.DefaultMemPool = &allocator{} // mempool.New(1024+1024, 1024*1024*1024)
	} else {
		mempool.DefaultMemPool = mempool.New(cnf.ReadBufferSize+1024, 1024*1024*1024)
	}

	upgrader.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		c.WriteMessage(messageType, data)
	})
	upgrader.KeepaliveTime = 0

	addrs, err := config.GetFrameworkServerAddrs(config.NbioModNonblocking, cnf.LimitPortRange)
	if err != nil {
		log.Fatalf("GetFrameworkBenchmarkAddrs(%v) failed: %v", config.NbioModNonblocking, err)
	}
	engine := cnf.startServers(addrs)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	engine.Shutdown(ctx)
}

func onWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}
	c.SetReadDeadline(time.Time{})
}

func (c *Config) startServers(addrs []string) *nbhttp.Engine {
	mux := &http.ServeMux{}
	mux.HandleFunc("/ws", onWebsocket)
	core.HandleCommon(mux)
	engine := nbhttp.NewEngine(nbhttp.Config{
		Network:                 "tcp",
		Addrs:                   addrs,
		Handler:                 mux,
		IOMod:                   nbhttp.IOModNonBlocking,
		ReleaseWebsocketPayload: true,
		Listen:                  core.Listen2(c.Reuse),
	})

	err := engine.Start()
	if err != nil {
		log.Fatalf("nbio.Start failed: %v", err)
	}

	return engine
}

type allocator struct{}

func (a *allocator) Malloc(size int) []byte {
	return make([]byte, size)
}

func (a *allocator) Realloc(buf []byte, size int) []byte {
	if size <= cap(buf) {
		return buf[:size]
	}
	return append(buf, make([]byte, size-cap(buf))...)
}

func (a *allocator) Append(buf []byte, more ...byte) []byte {
	return append(buf, more...)
}

func (a *allocator) AppendString(buf []byte, more string) []byte {
	return append(buf, more...)
}

func (a *allocator) Free(buf []byte) {
}
