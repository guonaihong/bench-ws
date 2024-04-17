package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/guonaihong/bench-ws/config"
	"github.com/guonaihong/bench-ws/core"
	"github.com/guonaihong/clop"
	"github.com/lxzan/gws"
)

func startServers(addrs []string, reuse bool) []net.Listener {
	lns := make([]net.Listener, 0, len(addrs))
	for _, addr := range addrs {
		server := gws.NewServer(new(Handler), &gws.ServerOption{})
		ln, err := core.Listen("tcp", addr, reuse)
		if err != nil {
			log.Fatalf("Listen failed: %v", err)
		}
		lns = append(lns, ln)
		go func() {
			log.Printf("server exit: %v", server.RunListener(ln))
		}()
	}
	return lns
}
func main() {
	h := &Handler{}
	clop.Bind(h)

	addrs, err := config.GetFrameworkServerAddrs(config.Gws, h.LimitPortRange)
	if err != nil {
		log.Fatalf("GetFrameworkBenchmarkAddrs(%v) failed: %v", config.Gws, err)
	}
	lns := startServers(addrs, h.Reuse)
	pidServerAddr, err := config.GetFrameworkPidServerAddrs(config.Gws)
	if err != nil {
		log.Fatalf("GetFrameworkPidServerAddrs(%v) failed: %v", config.Gws, err)
	}
	var pidLn net.Listener
	go func() {
		mux := &http.ServeMux{}
		core.HandleCommon(mux)
		ln, err := core.Listen("tcp", pidServerAddr, h.Reuse)
		if err != nil {
			log.Fatalf("Listen failed: %v", err)
		}
		pidLn = ln
		log.Printf("pid server exit: %v", http.Serve(ln, mux))
	}()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	for _, ln := range lns {
		ln.Close()
	}
	pidLn.Close()
}

type Handler struct {
	gws.BuiltinEventHandler
	// 是否异步写
	AsyncWrite bool   `clop:"short;long" usage:"async write"`
	Addr       string `clop:"long" usage:"websocket server address" default:":6666""`
	// 打开tcp nodealy
	OpenTcpDelay bool `clop:"short;long" usage:"tcp delay"`
	Reuse        bool `clop:"short;long" usage:"reuse port"`
	// 使用限制端口范围, 默认1， -1表示不限制
	LimitPortRange int `clop:"short;long" usage:"limit port range" default:"1"`
}

func (c *Handler) OnOpen(socket *gws.Conn) {
	if c.OpenTcpDelay {
		socket.SetNoDelay(!c.OpenTcpDelay)
	}
}

func (c *Handler) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.WritePong(payload)
}

func (c *Handler) OnMessage(socket *gws.Conn, message *gws.Message) {

	if c.AsyncWrite {
		socket.WriteAsync(message.Opcode, message.Bytes(), func(err error) {
			message.Close()
		})
	} else {
		_ = socket.WriteMessage(message.Opcode, message.Bytes())
	}
	message.Close()
}
