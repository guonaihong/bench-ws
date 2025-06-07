package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	nettyws "github.com/go-netty/go-netty-ws"
	"github.com/guonaihong/bench-ws/core"
	"github.com/guonaihong/bench-ws/pkg/port"
	"github.com/guonaihong/clop"
)

type Conf struct {
	Addr    string `clop:"short;long" usage:"websocket server address" default:":6666""`
	Nodelay bool   `clop:"short;long usage:"nodelay" default:"true"`
	core.BaseCmd
}

// https://github.com/go-netty/go-netty-ws/blob/master/example/echo.go
func main() {
	// serveMux := http.NewServeMux()

	// ws := nettyws.NewWebsocket(
	// 	nettyws.WithServeMux(serveMux),
	// 	nettyws.WithBinary(),
	// 	nettyws.WithBufferSize(4096, 0),
	// )

	cnf := &Conf{}
	clop.Bind(cnf)

	portRange, err := port.GetPortRange("NETTYWS")
	if err != nil {
		log.Fatalf("GetPortRange(%v) failed: %v", "NETTYWS", err)
	}

	wg := sync.WaitGroup{}
	defer wg.Wait()
	for port := portRange.Start; port <= portRange.End; port++ {
		wg.Add(1)
		go cnf.startServer(port, &wg)
	}
}

func (c *Conf) startServer(port int, wg *sync.WaitGroup) {

	defer wg.Done()

	mux := &http.ServeMux{}

	var ws = nettyws.NewWebsocket(
		nettyws.WithServeMux(mux),
		nettyws.WithBinary(),
		nettyws.WithNoDelay(c.Nodelay),
		nettyws.WithBufferSize(2048, 0),
	)

	ws.OnData = func(conn nettyws.Conn, data []byte) {
		conn.Write(data)
	}
	ws.Listen(fmt.Sprintf(":%d", port))
}
