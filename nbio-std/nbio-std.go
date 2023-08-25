package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/guonaihong/clop"
	"github.com/lesismal/nbio/logging"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

var upgrader = websocket.NewUpgrader()

type Config struct {
	Addr string `clop:"short;long" usage:"websocket server address" default:":4444""`
}

func main() {
	var conf Config
	clop.Bind(&conf)
	logging.SetLevel(logging.LevelError)
	upgrader.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		c.WriteMessage(messageType, data)
	})
	upgrader.BlockingModAsyncWrite = false

	mux := &http.ServeMux{}
	mux.HandleFunc("/", onWebsocket)

	rawTCP, err := net.Listen("tcp", conf.Addr)
	if err != nil {
		fmt.Println("Listen fail:", err)
		return
	}

	log.Println("non-tls server exit:", http.Serve(rawTCP, mux))
}

func onWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}
	c.SetReadDeadline(time.Time{})
}
