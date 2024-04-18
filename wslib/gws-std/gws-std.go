package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/guonaihong/bench-ws/config"
	"github.com/guonaihong/bench-ws/core"
	"github.com/guonaihong/clop"
	"github.com/lxzan/gws"
)

type Config struct {
	gws.BuiltinEventHandler
	// 是否异步写
	AsyncWrite bool   `clop:"short;long" usage:"async write"`
	Addr       string `clop:"long" usage:"websocket server address" default:":6666""`
	// 打开tcp nodealy
	OpenTcpDelay bool `clop:"short;long" usage:"tcp delay"`

	core.BaseCmd
	upgrader *gws.Upgrader
}

func (c *Config) OnOpen(socket *gws.Conn) {
	if c.OpenTcpDelay {
		socket.SetNoDelay(!c.OpenTcpDelay)
	}
}

func (c *Config) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.WritePong(payload)
}

func (c *Config) OnMessage(socket *gws.Conn, message *gws.Message) {

	if c.AsyncWrite {
		socket.WriteAsync(message.Opcode, message.Bytes(), func(err error) {
			message.Close()
		})
	} else {
		_ = socket.WriteMessage(message.Opcode, message.Bytes())
	}
	message.Close()
}

func (c *Config) onWebsocket(w http.ResponseWriter, r *http.Request) {
	con, err := c.upgrader.Upgrade(w, r)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}
	// frameworks.SetNoDelay(c.NetConn(), *nodelay)
	// c.SetReadDeadline(time.Time{})
	go func() {
		con.ReadLoop()
	}()
}
func main() {
	cnf := &Config{}
	clop.Bind(cnf)
	cnf.upgrader = gws.NewUpgrader(cnf, &gws.ServerOption{})

	addrs, err := config.GetFrameworkServerAddrs(config.GwsStd, cnf.LimitPortRange)
	if err != nil {
		log.Fatalf("GetFrameworkBenchmarkAddrs(%v) failed: %v", config.GwsStd, err)
	}
	lns := core.StartServers(addrs, cnf.onWebsocket, cnf.Reuse)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	for _, ln := range lns {
		ln.Close()
	}

}
