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
	AsyncWrite bool `clop:"short;long" usage:"async write"`
}

func (c *Handler) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.WritePong(payload)
}

func (c *Handler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	if c.AsyncWrite {
		_ = socket.WriteAsync(message.Opcode, message.Bytes())
	} else {
		_ = socket.WriteMessage(message.Opcode, message.Bytes())
	}
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

	if err := http.ListenAndServe(":6666", nil); err != nil {
		log.Fatalf("%v", err)
	}
}
