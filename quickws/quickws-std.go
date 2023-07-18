//go:build quickwstest
// +build quickwstest

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	// _ "net/http/pprof"

	"github.com/antlabs/quickws"
	"github.com/guonaihong/clop"
	//"os"
)

type Config struct {
	DisableUtf8 bool `clop:"short;long" usage:"disable utf8"`
	// 几倍的窗口大小
	WindowsMultipleTimesPayloadSize int `clop:"short;long" usage:"windows multiple times payload size"`
	// 使用bufio的解析方式
	UseBufio bool `clop:"short;long" usage:"use bufio"`
	// 打开tcp nodealy
	OpenDelay bool `clop:"short;long" usage:"tcp no delay"`
}

type echoHandler struct{}

func (e *echoHandler) OnOpen(c *quickws.Conn) {
	// fmt.Printf("OnOpen: %p\n", c)
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
func (cnf *Config) echo(w http.ResponseWriter, r *http.Request) {
	size := float32(1.0)
	if cnf.WindowsMultipleTimesPayloadSize > 0 {
		size = float32(cnf.WindowsMultipleTimesPayloadSize)
	}
	opt := []quickws.ServerOption{
		quickws.WithServerReplyPing(),
		// quickws.WithServerDecompression(),
		quickws.WithServerIgnorePong(),
		quickws.WithServerCallback(&echoHandler{}),
		// quickws.WithServerReadTimeout(5*time.Second),
		quickws.WithWindowsMultipleTimesPayloadSize(size),
	}

	if cnf.OpenDelay {
		opt = append(opt, quickws.WithServerTCPDelay())
	}
	if cnf.UseBufio {
		opt = append(opt, quickws.WithBufioParseMode())
	}

	if cnf.DisableUtf8 {
		opt = append(opt, quickws.WithServerDisableUTF8Check())
	}

	c, err := quickws.Upgrade(w, r, opt...)
	if err != nil {
		fmt.Println("Upgrade fail:", err)
		return
	}

	c.StartReadLoop()
}

func main() {
	var conf Config
	clop.Bind(&conf)

	mux := &http.ServeMux{}
	mux.HandleFunc("/", conf.echo)

	go func() {
		// log.Println(http.ListenAndServe(":6060", nil))
	}()
	rawTCP, err := net.Listen("tcp", ":9001")
	if err != nil {
		fmt.Println("Listen fail:", err)
		return
	}

	log.Println("non-tls server exit:", http.Serve(rawTCP, mux))
}
