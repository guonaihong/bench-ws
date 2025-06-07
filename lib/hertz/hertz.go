package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/guonaihong/bench-ws/core"
	"github.com/guonaihong/bench-ws/pkg/port"
	"github.com/guonaihong/clop"
	"github.com/hertz-contrib/websocket"
)

var upgrader = websocket.HertzUpgrader{}

type Config struct {
	// 打开性能优化开关
	UseReader bool   `clop:"short;long" usage:"use reader"`
	Addr      string `clop:"short;long" usage:"websocket server address" default:":5555"`
	// 打开tcp nodealy
	OpenTcpDelay bool `clop:"short;long" usage:"tcp delay"`
	core.BaseCmd
}

func setNoDelay(c net.Conn, noDelay bool) error {
	if tcp, ok := c.(*net.TCPConn); ok {
		return tcp.SetNoDelay(noDelay)
	}
	return nil
}

func (c *Config) work(conn *websocket.Conn) {
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
			if err != nil {
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
}

func (c *Config) startServer(port int, wg *sync.WaitGroup) {
	defer wg.Done()

	h := server.Default(
		server.WithHostPorts(fmt.Sprintf(":%d", port)),
	)

	h.GET("/", func(ctx context.Context, req *app.RequestContext) {
		err := upgrader.Upgrade(req, func(conn *websocket.Conn) {
			setNoDelay(conn.NetConn(), !c.OpenTcpDelay)
			c.work(conn)
		})
		if err != nil {
			log.Printf("upgrade failed: %v", err)
		}
	})

	go func() {
		h.Spin()
	}()
}

func main() {
	var cnf Config
	clop.Bind(&cnf)

	portRange, err := port.GetPortRange("HERTZ")
	if err != nil {
		log.Fatalf("GetPortRange(%v) failed: %v", "HERTZ", err)
	}

	wg := sync.WaitGroup{}
	for port := portRange.Start; port <= portRange.End; port++ {
		wg.Add(1)
		go cnf.startServer(port, &wg)
	}
	wg.Wait()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
