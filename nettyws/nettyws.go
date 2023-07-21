package main

import (
	"log"
	"net/http"

	nettyws "github.com/go-netty/go-netty-ws"
)

func main() {
	serveMux := http.NewServeMux()

	ws := nettyws.NewWebsocket(
		nettyws.WithServeMux(serveMux),
		nettyws.WithBinary(),
		nettyws.WithBufferSize(2048, 2048),
	)
	ws.OnData = func(conn nettyws.Conn, data []byte) {
		conn.Write(data)
	}
	// addr := fmt.Sprintf("%s/ws", addr)
	log.Printf("server exit: %v", ws.Listen(":5001"))
}
