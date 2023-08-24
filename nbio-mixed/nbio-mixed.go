package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/guonaihong/clop"
	"github.com/lesismal/nbio/mempool"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

var upgrader = websocket.NewUpgrader()

type Config struct {
	// 打开性能优化开关
	UseReader      bool `clop:"short;long" usage:"use reader"`
	ReadBufferSize int  `clop:"short;long" usage:"read buffer size" default:"1024"`

	Addr string `clop:"short;long" usage:"websocket server address" default:":4444""`
}

func main() {
	var conf Config
	clop.Bind(&conf)
	mempool.DefaultMemPool = mempool.New(1024+1024, 1024*1024*1024)

	upgrader.BlockingModAsyncWrite = false

	upgrader.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		c.WriteMessage(messageType, data)
	})

	mux := &http.ServeMux{}
	mux.HandleFunc("/", onWebsocket)

	engine := nbhttp.NewEngine(nbhttp.Config{
		Network: "tcp",
		Addrs:   []string{conf.Addr},
		Handler: mux,
		// IOMod:                   nbhttp.IOModBlocking,
		IOMod:                   nbhttp.IOModMixed,
		ReleaseWebsocketPayload: true,
		Listen:                  net.Listen,
	})
	engine.Start()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	engine.Stop()
}

func onWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}
	c.SetReadDeadline(time.Time{})
}
