package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	nettyws "github.com/go-netty/go-netty-ws"
	"github.com/guonaihong/bench-ws/config"
	"github.com/guonaihong/bench-ws/core"
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

	addrs, err := config.GetFrameworkServerAddrs(config.GoNettyWs, cnf.LimitPortRange)
	if err != nil {
		log.Fatalf("GetFrameworkBenchmarkAddrs(%v) failed: %v", config.GoNettyWs, err)
	}
	svrs := cnf.startServers(addrs)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	for _, svr := range svrs {
		svr.Close()
	}

}

func (c *Conf) startServers(addrs []string) []*nettyws.Websocket {
	svrs := make([]*nettyws.Websocket, 0, len(addrs))
	for _, addr := range addrs {
		var mux = http.NewServeMux()
		core.HandleCommon(mux)
		var ws = nettyws.NewWebsocket(
			nettyws.WithServeMux(mux),
			nettyws.WithBinary(),
			nettyws.WithNoDelay(c.Nodelay),
			nettyws.WithBufferSize(2048, 0),
		)
		svrs = append(svrs, ws)
		ws.OnData = func(conn nettyws.Conn, data []byte) {
			conn.Write(data)
		}
		addr := fmt.Sprintf("%s/ws", addr)
		go func() {
			log.Printf("server exit: %v", ws.Listen(addr))
		}()
	}
	return svrs
}
