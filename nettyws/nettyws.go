package main

import (
	"log"

	nettyws "github.com/go-netty/go-netty-ws"
	"github.com/guonaihong/clop"
)

type Conf struct {
	Addr string `clop:"short;long" usage:"websocket server address" default:":6666""`
}

// https://github.com/go-netty/go-netty-ws/blob/master/example/echo.go
func main() {
	// serveMux := http.NewServeMux()

	// ws := nettyws.NewWebsocket(
	// 	nettyws.WithServeMux(serveMux),
	// 	nettyws.WithBinary(),
	// 	nettyws.WithBufferSize(4096, 0),
	// )

	c := &Conf{}
	clop.Bind(c)

	ws := nettyws.NewWebsocket(
		nettyws.WithBinary(),
		nettyws.WithBufferSize(4096, 0),
	)
	ws.OnData = func(conn nettyws.Conn, data []byte) {
		conn.Write(data)
	}
	// addr := fmt.Sprintf("%s/ws", addr)
	log.Printf("server exit: %v", ws.Listen(c.Addr))
}
