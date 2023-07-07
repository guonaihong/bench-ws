//go:build quickwstest
// +build quickwstest

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/antlabs/quickws"
	//"os"
)

type echoHandler struct{}

func (e *echoHandler) OnOpen(c *quickws.Conn) {
	fmt.Printf("OnOpen: %p\n", c)
}

func (e *echoHandler) OnMessage(c *quickws.Conn, op quickws.Opcode, msg []byte) {
	// fmt.Println("OnMessage:", msg)
	// if err := c.WriteTimeout(op, msg, 3*time.Second); err != nil {
	// 	fmt.Println("write fail:", err)
	// }
	c.WriteMessage(op, msg)
}

func (e *echoHandler) OnClose(c *quickws.Conn, err error) {
	fmt.Printf("OnClose:%p, %v\n", c, err)
}

// echo测试服务
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := quickws.Upgrade(w, r,
		quickws.WithServerReplyPing(),
		quickws.WithServerDecompression(),
		quickws.WithServerIgnorePong(),
		quickws.WithServerCallback(&echoHandler{}),
		quickws.WithServerReadTimeout(5*time.Second),
	)
	if err != nil {
		fmt.Println("Upgrade fail:", err)
		return
	}

	c.ReadLoop()
}

func main() {
	mux := &http.ServeMux{}
	mux.HandleFunc("/autobahn", echo)

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	rawTCP, err := net.Listen("tcp", ":9001")
	if err != nil {
		fmt.Println("Listen fail:", err)
		return
	}

	log.Println("non-tls server exit:", http.Serve(rawTCP, mux))
}
