package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/bytedance/gopkg/util/gopool"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/guonaihong/clop"
	"github.com/hertz-contrib/websocket"
)

var upgrader = websocket.HertzUpgrader{}

type Config struct {
	// 打开性能优化开关
	UseReader      bool `clop:"short;long" usage:"use reader"`
	ReadBufferSize int  `clop:"short;long" usage:"read buffer size" default:"1024"`

	Addr string `clop:"short;long" usage:"websocket server address" default:":5555""`
}

func main() {
	var conf Config
	clop.Bind(&conf)

	gopool.SetCap(1000000)

	srv := server.New(server.WithHostPorts(conf.Addr))
	go func() {
		srv.GET("/", conf.onWebsocket)
		srv.Spin()
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
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
