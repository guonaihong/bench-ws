package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/guonaihong/bench-ws/core"
	"github.com/guonaihong/bench-ws/pkg/port"
	"github.com/guonaihong/clop"
	"github.com/lesismal/nbio/mempool"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

var upgrader = websocket.NewUpgrader()

type Config struct {
	Addr string `clop:"short;long" usage:"websocket server address" default:":4444""`
	core.BaseCmd
}

func main() {
	var cnf Config
	clop.Bind(&cnf)

	mempool.DefaultMemPool = mempool.New(cnf.ReadBufferSize+1024, 1024*1024*1024)

	upgrader.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		c.WriteMessage(messageType, data)
	})
	upgrader.KeepaliveTime = 0
	upgrader.BlockingModAsyncWrite = false

	portRange, err := port.GetPortRange("NBIO-BLOCKING")
	if err != nil {
		log.Fatalf("GetPortRange(%v) failed: %v", "NBIO-BLOCKING", err)
	}
	fmt.Println("NBIO-BLOCKING server starting on ports", portRange.Start, "-", portRange.End)
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

func (c *Config) startServer(port int, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := &http.ServeMux{}
	mux.HandleFunc("/", onWebsocket)
	engine := nbhttp.NewEngine(nbhttp.Config{
		Network:                 "tcp",
		Addrs:                   []string{fmt.Sprintf(":%d", port)},
		Handler:                 mux,
		IOMod:                   nbhttp.IOModBlocking,
		ReleaseWebsocketPayload: true,
		// Listen:                  core.Listen2(c.Reuse),
	})

	err := engine.Start()
	if err != nil {
		log.Fatalf("nbio.Start failed: %v", err)
	}

}
func onWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}
	c.SetReadDeadline(time.Time{})
}
