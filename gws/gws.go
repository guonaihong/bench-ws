package main

import (
	"log"

	"github.com/guonaihong/clop"
	"github.com/lxzan/gws"
)

func main() {
	h := &Handler{}
	clop.Bind(h)
	app := gws.NewServer(h, &gws.ServerOption{
		// CompressEnabled:  true,
		CheckUtf8Enabled: false,
	})
	log.Fatalf("%v", app.Run(":8001"))
}

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
