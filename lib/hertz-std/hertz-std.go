package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/bytedance/gopkg/util/gopool"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/guonaihong/bench-ws/core"
	"github.com/guonaihong/bench-ws/pkg/port"
	"github.com/guonaihong/clop"
	"github.com/hertz-contrib/pprof"
	"github.com/hertz-contrib/websocket"
)

var upgrader = websocket.HertzUpgrader{}

type Config struct {
	// 打开性能优化开关
	UseReader bool   `clop:"short;long" usage:"use reader"`
	Addr      string `clop:"short;long" usage:"websocket server address" default:":5555""`
	core.BaseCmd
}

func main() {
	var cnf Config
	clop.Bind(&cnf)

	hlog.SetLevel(hlog.LevelFatal)
	gopool.SetCap(1000000)

	portRange, err := port.GetPortRange("HERTZ-STD")
	if err != nil {
		log.Fatalf("GetPortRange(%v) failed: %v", "HERTZ-STD", err)
	}
	fmt.Println("HERTZ-STD server starting on ports", portRange.Start, "-", portRange.End)
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

func onServerPid(c context.Context, ctx *app.RequestContext) {
	ctx.Response.BodyWriter().Write([]byte(fmt.Sprintf("%d", os.Getpid())))
}

func (c *Config) startServer(port int, wg *sync.WaitGroup) {
	defer wg.Done()

	srv := server.New(server.WithHostPorts(fmt.Sprintf(":%d", port)),
		server.WithTransport(standard.NewTransporter))
	pprof.Register(srv)
	srv.GET("/", c.onWebsocket)
	srv.GET("/pid", onServerPid)
	srv.Spin()
}

func (c *Config) onWebsocket(c2 context.Context, ctx *app.RequestContext) {
	upgradeErr := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		defer conn.Close()

		if !c.UseReader {
			for {
				mt, message, err := conn.ReadMessage()
				if err != nil {
					// log.Printf("read message failed: %v", err)
					return
				}
				err = conn.WriteMessage(mt, message)
				if err != nil {
					// log.Printf("write failed: %v", err)
					return
				}
			}
		}
		var nread int
		buffer := make([]byte, c.ReadBufferSize)
		readBuffer := buffer
		for {
			mt, reader, err := conn.NextReader()
			if err != nil {
				// log.Printf("read failed: %v", err)
				return
			}
			for {
				if nread == len(readBuffer) {
					readBuffer = append(readBuffer, buffer...)
				}
				n, err := reader.Read(readBuffer[nread:])
				nread += n
				if err == io.EOF {
					break
				}
			}
			err = conn.WriteMessage(mt, readBuffer[:nread])
			nread = 0
			if err != nil {
				// log.Printf("write failed: %v", err)
				return
			}
		}
	})

	if upgradeErr != nil {
		log.Printf("upgrade failed: %v", upgradeErr)
		return
	}
}
