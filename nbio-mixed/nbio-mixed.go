package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/guonaihong/bench-ws/config"
	"github.com/guonaihong/bench-ws/core"
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

	Addr              string `clop:"short;long" usage:"websocket server address" default:":4444""`
	MaxBlockingOnline int    `clop:"short;long" usage:"max blocking online num, e.g. 10000" default:"10000"`
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

	addrs, err := config.GetFrameworkServerAddrs(config.NbioModMixed, cnf.LimitPortRange)
	if err != nil {
		log.Fatalf("GetFrameworkBenchmarkAddrs(%v) failed: %v", config.NbioModMixed, err)
	}
	engine := cnf.startServers(addrs)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	engine.Stop()
}

func (c *Config) startServers(addrs []string) *nbhttp.Engine {
	mux := &http.ServeMux{}
	mux.HandleFunc("/ws", onWebsocket)
	core.HandleCommon(mux)
	engine := nbhttp.NewEngine(nbhttp.Config{
		Network:                 "tcp",
		Addrs:                   addrs,
		Handler:                 mux,
		IOMod:                   nbhttp.IOModMixed,
		MaxBlockingOnline:       c.MaxBlockingOnline,
		ReleaseWebsocketPayload: true,
		Listen:                  core.Listen2(c.Reuse),
	})

	err := engine.Start()
	if err != nil {
		log.Fatalf("nbio.Start failed: %v", err)
	}

	return engine
}
func onWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}
	c.SetReadDeadline(time.Time{})
}
