package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"time"

	"github.com/guonaihong/bench-ws/core"
	"github.com/guonaihong/bench-ws/pkg/port"
	"github.com/guonaihong/clop"

	// "github.com/lesismal/nbio/log"
	"github.com/lesismal/nbio/mempool"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

var upgrader = websocket.NewUpgrader()

type Config struct {
	// 打开性能优化开关
	UseStdMalloc bool   `clop:"short;long" usage:"use reader"`
	Addr         string `clop:"short;long" usage:"websocket server address" default:":4444""`
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

	portRange, err := port.GetPortRange("NBIO-NONBLOCKING")
	if err != nil {
		log.Fatalf("GetPortRange(%v) failed: %v", "NBIO-NONBLOCKING", err)
	}
	fmt.Println("NBIO-NONBLOCKING server starting on ports", portRange.Start, "-", portRange.End)

	wg := sync.WaitGroup{}
	defer wg.Wait()
	for port := portRange.Start; port <= portRange.End; port++ {
		wg.Add(1)
		go cnf.startServer(port, &wg)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}

func onWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}
	c.SetReadDeadline(time.Time{})
}

func (c *Config) startServer(port int, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := &http.ServeMux{}
	mux.HandleFunc("/", onWebsocket)
	engine := nbhttp.NewEngine(nbhttp.Config{
		Network:                 "tcp",
		Addrs:                   []string{fmt.Sprintf(":%d", port)},
		Handler:                 mux,
		IOMod:                   nbhttp.IOModNonBlocking,
		ReleaseWebsocketPayload: true,
		// Listen:                  core.Listen2(c.Reuse),
	})

	err := engine.Start()
	if err != nil {
		log.Fatalf("nbio.Start failed: %v", err)
	}

}

type allocator struct{}

func (a *allocator) Malloc(size int) *[]byte {
	rv := make([]byte, size)
	return &rv
}

func (a *allocator) Realloc(buf *[]byte, size int) *[]byte {
	if size <= cap(*buf) {
		rv := (*buf)[:size]
		return &rv
	}

	rv := append(*buf, make([]byte, size-cap(*buf))...)
	return &rv
}

func (a *allocator) Append(buf *[]byte, more ...byte) *[]byte {
	rv := append(*buf, more...)
	return &rv
}

func (a *allocator) AppendString(buf *[]byte, more string) *[]byte {
	rv := append(*buf, more...)
	return &rv
}

func (a *allocator) Free(buf *[]byte) {
}
