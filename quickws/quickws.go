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
	EnableUtf8 bool `clop:"short;long" usage:"enable utf8"`
	// 几倍的窗口大小
	WindowsMultipleTimesPayloadSize int `clop:"short;long" usage:"windows multiple times payload size"`
	// 使用bufio的解析方式
	UseBufio bool `clop:"short;long" usage:"use bufio"`
	// 打开tcp nodealy
	OpenDelay bool `clop:"short;long" usage:"tcp delay"`
	// 关闭bufio clear hack 优化
	DisableBufioClearHack bool `clop:"long" usage:"disable bufio clear hack"`
}

var upgrader *quickws.UpgradeServer

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
	c, err := upgrader.Upgrade(w, r)
	if err != nil {
		fmt.Println("Upgrade fail:", err)
		return
	}

	c.StartReadLoop()
}

func main() {
	var cnf Config
	clop.Bind(&cnf)

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
		quickws.WithServerWindowsMultipleTimesPayloadSize(size),
	}

	if cnf.OpenDelay {
		opt = append(opt, quickws.WithServerTCPDelay())
	}

	if cnf.UseBufio {
		opt = append(opt, quickws.WithServerBufioParseMode())
	}

	if cnf.EnableUtf8 {
		opt = append(opt, quickws.WithServerEnableUTF8Check())
	}

	if cnf.DisableBufioClearHack {
		opt = append(opt, quickws.WithServerDisableBufioClearHack())
	}

	upgrader = quickws.NewUpgrade(opt...)

	mux := &http.ServeMux{}
	mux.HandleFunc("/", cnf.echo)

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
