package main

import (
	"log"
	"net/http"

	"github.com/guonaihong/clop"
	"github.com/lxzan/gws"
)

type Handler struct {
	gws.BuiltinEventHandler
	// 是否异步写
	AsyncWrite bool   `clop:"short;long" usage:"async write"`
	Addr       string `clop:"long" usage:"websocket server address" default:":6666""`
	// 打开tcp nodealy
	OpenTcpDelay bool `clop:"short;long" usage:"tcp delay"`
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

func main() {
	h := &Handler{}
	clop.Bind(h)
	upgrader := gws.NewUpgrader(h, &gws.ServerOption{})

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		socket, err := upgrader.Upgrade(writer, request)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		socket.ReadLoop()
	})

	if err := http.ListenAndServe(h.Addr, nil); err != nil {
		log.Fatalf("%v", err)
	}
}
